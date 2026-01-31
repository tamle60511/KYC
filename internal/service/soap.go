package service

import (
	"CQS-KYC/config"
	"CQS-KYC/database"
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/clbanning/mxj/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ERPService struct {
	db             database.Database
	config         *config.Config
	userRepo       repository.UserRepo
	wfDefSerivce   WorkflowService
	workflowEngine InstanceService
}

func NewERPService(
	db database.Database,
	config *config.Config,
	userRepo repository.UserRepo,
	wfDefSerivce WorkflowService,
	workflowEngine InstanceService,
) *ERPService {
	return &ERPService{
		db:             db,
		config:         config,
		userRepo:       userRepo,
		wfDefSerivce:   wfDefSerivce,
		workflowEngine: workflowEngine,
	}
}

type ResultItem struct {
	Key, Value, Source string
}

type ExtractedData struct {
	CompanyId   string
	FormId      string
	ComPRID     string
	UserID      string
	DocType     string
	DocNum      string
	Action      string
	SSLProtocal string
	RawData     map[string]interface{}
}

// =============================================================================
// 1. MAIN FLOW: NHẬN REQUEST -> BÓC TÁCH -> CHẠY ENGINE -> TRẢ LỜI ERP
// =============================================================================
func (s *ERPService) ProcessSOAPRequest(xmlBody []byte) error {
	// 1. Bóc tách dữ liệu
	data, err := s.processAndExtract(xmlBody)
	if err != nil {
		fmt.Printf("[PARSE ERROR] %v\n", err)
		return nil // Return nil để SOAP Handler trả về 200 OK (tránh ERP gửi lại nếu lỗi format)
	}

	fmt.Printf("[RECEIVED] %s | Type: %s | Num: %s | User: %s\n", data.CompanyId, data.DocType, data.DocNum, data.UserID)

	// 2. Kích hoạt Workflow Engine
	if err := s.routeAndInitiateWorkflow(data); err != nil {
		fmt.Printf("[WORKFLOW FAIL] %v\n", err)
		// Có thể return err để ERP biết lỗi, hoặc return nil và log lại tùy nghiệp vụ
		return err
	}

	// 3. Quan trọng: Báo cho ERP biết đã nhận đơn thành công (Chặn retry)
	// Hàm này update vào bảng Queue của ERP (EFJobQue/DSCSYS)
	go func() {
		if err := s.updateStatusInDSCSYS(data); err != nil {
			fmt.Printf(" [WARN] Failed to update EFJobQue: %v\n", err)
		}
	}()

	fmt.Println(" [DONE] Process completed successfully.")
	return nil
}

// =============================================================================
// 2. CORE LOGIC: ROUTING & INITIATION
// =============================================================================
func (s *ERPService) routeAndInitiateWorkflow(data *ExtractedData) error {
	ctx := context.Background()
	jsonBytes, err := json.Marshal(data.RawData)
	if err != nil {
		return fmt.Errorf("marshal json failed: %w", err)
	}

	return s.db.DB().Transaction(func(tx *gorm.DB) error {
		// ---------------------------------------------------------
		// BƯỚC 1: KIỂM TRA & KHÓA BẢN GHI (BLOCKING DUPLICATE)
		// ---------------------------------------------------------
		var req model.Request

		// Tìm xem request đã có chưa
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("doc_num = ? AND doc_type = ?", data.DocNum, data.DocType).
			First(&req).Error

		if err == nil {
			// A. Nếu tìm thấy record
			if req.WorkflowInstanceID != 0 {
				fmt.Printf("[DUPLICATE] Request %s existed with InstanceID: %d. Ignored.\n", data.DocNum, req.WorkflowInstanceID)
				return nil // Trả về nil để báo thành công giả, chặn ERP retry
			}
			// Nếu record tồn tại nhưng chưa có InstanceID (có thể do lần trước crash giữa chừng),
			// ta sẽ dùng lại record 'req' này để xử lý tiếp.
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Lỗi DB thật
			return err
		}

		// ---------------------------------------------------------
		// BƯỚC 2: MAP USER (STRICT MODE)
		// ---------------------------------------------------------
		var creatorIDStr string
		var deptID uint64
		var factoryID uint64

		user, err := s.userRepo.GetByCode(ctx, data.UserID)
		if err != nil {
			// OPTION A: Hard Fail (Khuyên dùng) - Bắt buộc data phải chuẩn
			// return fmt.Errorf("user %s not found in system - sync required", data.UserID)

			// OPTION B: Fallback an toàn (Gán cho Admin nhưng Log cảnh báo rõ)
			fmt.Printf("[WARN] User %s not found. Assigning to Admin (ID=1) for manual check.\n", data.UserID)
			creatorIDStr = "1"
			deptID = 1    // Default Dept của Admin
			factoryID = 1 // Default Factory của Admin
		} else {
			creatorIDStr = fmt.Sprintf("%d", user.ID)
			deptID = user.DepartmentID
			factoryID = user.FactoryID
		}

		// ---------------------------------------------------------
		// BƯỚC 3: TẠO/UPDATE REQUEST (STAGING)
		// ---------------------------------------------------------
		// Nếu chưa có record (err == recordNotFound ở trên), tạo mới
		if req.ID == 0 {
			req = model.Request{
				ServiceName: data.FormId,
				CompanyID:   data.CompanyId,
				Operation:   data.ComPRID,
				DocType:     data.DocType,
				DocNum:      data.DocNum,
				CreatorID:   data.UserID,
				Status:      "INITIATING",
				SSLProtocal: data.SSLProtocal,
				Detail:      datatypes.JSON(jsonBytes),
			}
			if err := tx.Create(&req).Error; err != nil {
				return fmt.Errorf("create request failed: %w", err)
			}
		}

		// ---------------------------------------------------------
		// BƯỚC 4: LẤY DEFINITION & KHỞI TẠO WORKFLOW
		// ---------------------------------------------------------
		wfDef, err := s.wfDefSerivce.GetByCode(ctx, data.FormId)
		if err != nil {
			return fmt.Errorf("workflow definition not found for FormID: %s", data.FormId)
		}

		instance, err := s.workflowEngine.InitiateWorkflow(
			tx, // Pass transaction vào engine
			wfDef.ID,
			data.FormId,
			data.DocNum,
			data.DocType,
			creatorIDStr,
			factoryID,
			deptID,
			jsonBytes,
			"ERP_SOAP",
			"ERP_SYSTEM",
		)
		if err != nil {
			return fmt.Errorf("engine initiate failed: %w", err)
		}

		// ---------------------------------------------------------
		// BƯỚC 5: UPDATE LIÊN KẾT REQUEST -> INSTANCE
		// ---------------------------------------------------------
		req.WorkflowInstanceID = instance.ID
		req.Status = "PROCESSING"
		// Chỉ update các field cần thiết
		if err := tx.Model(&req).Updates(map[string]interface{}{
			"workflow_instance_id": instance.ID,
			"status":               "PROCESSING",
		}).Error; err != nil {
			return err
		}

		// ---------------------------------------------------------
		// BƯỚC 6: KHÓA ĐƠN TRÊN ERP (DB TRỰC TIẾP)
		// ---------------------------------------------------------
		if err := s.updateStatusInDB(data); err != nil {
			return fmt.Errorf("sync business status failed: %w", err)
		}

		return nil
	})
}

// =============================================================================
// 3. HELPERS: DB UPDATES
// =============================================================================

// Update bảng EFJobQue (DSCSYS) để báo ERP là "Tao nhận xong rồi, đừng gửi nữa"
func (s *ERPService) updateStatusInDSCSYS(req *ExtractedData) error {
	// Logic ghép chuỗi khóa chính của EFNET: DocType + "||" + DocNum (Thường là vậy)
	// Tuy nhiên, mẫu JSON của em có key "Condition" nhưng không thấy dùng.
	// Dựa vào logic cũ:
	efcondition := fmt.Sprintf("%s||%s", req.DocType, req.DocNum)

	updates := dto.EFJobQue{
		EF006: "Y", // Y = Success
		EF007: "Received by GO-Workflow",
	}

	// Lưu ý: Cần chắc chắn EF001 map đúng với CompanyId ("TESTEFNET")
	result := s.db.ERPDB().Model(&dto.EFJobQue{}).
		Where("EF001 = ? AND EF003 = ?", req.CompanyId, efcondition).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// Log warning thôi, không error vì có thể ERP đã xóa job rồi
		fmt.Printf("No EFJobQue record found to update for %s\n", efcondition)
	}
	return nil
}

// Update bảng nghiệp vụ (PURTA, COPTC...) set Trạng thái ký (TA016) = 1 (Đang trình ký)
func (s *ERPService) updateStatusInDB(data *ExtractedData) error {
	tableName, dbName, found := s.getTableAndDatabaseForDocType(data.ComPRID, data.CompanyId)
	if !found {
		return fmt.Errorf("config not found for ComPRID: %s, Company: %s", data.ComPRID, data.CompanyId)
	}

	updates := make(map[string]interface{})
	switch data.ComPRID {
	case "PURI05":
		updates["TA016"] = "1" // 1: Đang trình ký
	case "COPI06":
		updates["TC016"] = "1"
	default:
		// Mặc định cứ set TA016 nếu không rõ
		updates["TA016"] = "1"
	}

	// Construct table name: DBName.dbo.TableName
	fullTableName := fmt.Sprintf("%s.dbo.%s", dbName, tableName)

	err := s.db.ERPDB().Table(fullTableName).
		Where("TA001 = ? AND TA002 = ?", data.DocType, data.DocNum).
		Updates(updates).Error

	return err
}

func (s *ERPService) getTableAndDatabaseForDocType(comPRID, companyID string) (string, string, bool) {
	// 1. Map Program ID -> Table Name
	mappings := map[string]string{
		"PURI05": "PURTA", // Đơn mua hàng
		"COPI06": "COPTC", // Đơn hàng bán
	}
	tableName, ok := mappings[comPRID]
	if !ok {
		return "", "", false
	}

	// 2. Map Company ID -> Database Name
	// VÍ DỤ: "TESTEFNET" -> "TESTDB", "VN01" -> "ERPVN"
	// Em cần đảm bảo config.ERPDBMapping có chứa key lowercase của companyID
	companyLower := strings.ToLower(companyID)
	dbName, dbOk := s.config.ERPDBMapping[companyLower]

	if !dbOk {
		// Fallback: Nếu không tìm thấy mapping, thử dùng chính companyID làm tên DB
		// Hoặc trả về false tùy logic dự án
		fmt.Printf("DB Mapping not found for %s. Using default.\n", companyID)
		return tableName, companyID, true
	}

	return tableName, dbName, true
}

// =============================================================================
// 4. PARSING LOGIC (GIỮ NGUYÊN)
// =============================================================================
// ... (Giữ nguyên toàn bộ các hàm processAndExtract, recursiveScan, decodeAndFind, etc. từ code cũ của em)
// Vì phần này em làm tốt rồi, không cần sửa gì cả.

func (s *ERPService) processAndExtract(xmlData []byte) (*ExtractedData, error) {
	// ... (Copy lại y nguyên logic parsing cũ)
	mv, err := mxj.NewMapXml(xmlData)
	if err != nil {
		return nil, err
	}

	// Code tìm pParaContent ...
	paths := []string{
		"Envelope.Body.invokeSrv.pPara.#text",
		"Envelope.Body.invokeSrv.pPara",
		"Envelope.Body.InvokeSrv.pPara.#text",
		"Envelope.Body.InvokeSrv.pPara",
	}
	pParaContent := ""
	for _, p := range paths {
		if val, err := mv.ValueForPathString(p); err == nil && val != "" {
			pParaContent = val
			break
		}
	}
	if pParaContent == "" {
		return nil, fmt.Errorf("pPara empty")
	}

	results := make(map[string]ResultItem)
	// Decode logic...
	innerMap, err := mxj.NewMapXml([]byte(pParaContent))
	if err == nil {
		s.recursiveScan(map[string]interface{}(innerMap), results)
	} else {
		items := s.decodeAndFind(pParaContent)
		for _, item := range items {
			results[item.Key] = item
		}
	}

	// Populate Struct...
	finalData := &ExtractedData{RawData: make(map[string]interface{})}
	for k, v := range results {
		finalData.RawData[k] = v.Value
	}

	if v, ok := results["CompanyId"]; ok {
		finalData.CompanyId = v.Value
	}
	if v, ok := results["FormId"]; ok {
		finalData.FormId = v.Value
	}
	if v, ok := results["ComPRID"]; ok {
		finalData.ComPRID = v.Value
	}
	if v, ok := results["UserId"]; ok {
		finalData.UserID = v.Value
	}
	if v, ok := results["WhereClause"]; ok {
		finalData.DocType, finalData.DocNum = s.parseSQLWhere(v.Value)
	}

	// Fallback DocNum check
	if finalData.DocNum == "" {
		return nil, fmt.Errorf("missing DocNum")
	}

	return finalData, nil
}

// ... Copy nốt các hàm helper parseSQLWhere, decodeAndFind, flatScan, recursiveScan ...
// Chú ý: Nhớ copy đủ các hàm helper ở cuối file cũ sang đây nhé.
func (s *ERPService) recursiveScan(data interface{}, results map[string]ResultItem) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if strVal, ok := val.(string); ok && len(strVal) > 20 {
				items := s.decodeAndFind(strVal)
				for _, item := range items {
					results[item.Key] = item
				}
				if s.isKeyImportant(k) {
					results[k] = ResultItem{Key: k, Value: strVal, Source: "DirectMap"}
				}
			}
			s.recursiveScan(val, results)
		}
	case []interface{}:
		for _, val := range v {
			s.recursiveScan(val, results)
		}
	}
}

func (s *ERPService) decodeAndFind(content string) []ResultItem {
	var extracted []ResultItem
	content = strings.ReplaceAll(content, "]]>", "")
	content = strings.TrimSpace(content)

	if idx := strings.Index(content, " "); idx != -1 {
		parts := strings.Fields(content)
		longest := ""
		for _, p := range parts {
			if len(p) > len(longest) {
				longest = p
			}
		}
		content = longest
	}

	b, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		b, err = base64.URLEncoding.DecodeString(content)
		if err != nil {
			return extracted
		}
	}

	xmlStr := s.extractXMLString(string(b))
	if xmlStr == "" {
		return extracted
	}

	mv, err := mxj.NewMapXml([]byte(xmlStr))
	if err != nil {
		return extracted
	}

	s.flatScan(map[string]interface{}(mv), &extracted)
	return extracted
}

func (s *ERPService) flatScan(data interface{}, list *[]ResultItem) {
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			if strVal, ok := v.(string); ok {
				*list = append(*list, ResultItem{Key: k, Value: strVal, Source: "DecodedXML"})
			} else {
				s.flatScan(v, list)
			}
		}
	}
}

func (s *ERPService) extractXMLString(str string) string {
	start := strings.Index(str, "<")
	end := strings.LastIndex(str, ">")
	if start != -1 && end != -1 && end > start {
		return str[start : end+1]
	}
	return ""
}

func (s *ERPService) isKeyImportant(k string) bool {
	k = strings.ToLower(k)
	return strings.Contains(k, "company") || strings.Contains(k, "user") || strings.Contains(k, "where")
}

func (s *ERPService) parseSQLWhere(clause string) (string, string) {
	docTypePatterns := []string{
		`TA001\s*=\s*'([^']*)'`,
		`COPTC\.TC001\s*=\s*'([^']*)'`,
	}
	docNumPatterns := []string{
		`TA002\s*=\s*'([^']*)'`,
		`COPTC\.TC002\s*=\s*'([^']*)'`,
	}

	var docType, docNum string

	for _, pattern := range docTypePatterns {
		if matches := regexp.MustCompile(pattern).FindStringSubmatch(clause); len(matches) > 1 {
			docType = matches[1]
			break
		}
	}

	for _, pattern := range docNumPatterns {
		if matches := regexp.MustCompile(pattern).FindStringSubmatch(clause); len(matches) > 1 {
			docNum = matches[1]
			break
		}
	}

	return docType, docNum
}
