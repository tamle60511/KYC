package repository

import (
	"CQS-KYC/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type SignatureHelper struct {
	SecretKey *config.SignatureKeyConfig
}

func NewSignatureHelper(secretKey *config.SignatureKeyConfig) *SignatureHelper {
	return &SignatureHelper{SecretKey: secretKey}
}

func (s *SignatureHelper) GenerateSignature(
	signerID string,
	docNum string,
	action string,
	stepOrder int,
	requestData []byte, // Dữ liệu JSON của đơn
) (string, string, int64) {

	// 1. Hash dữ liệu (Snapshot)
	dataHash := ""
	if len(requestData) > 0 {
		sum := sha256.Sum256(requestData)
		dataHash = hex.EncodeToString(sum[:])
	}

	// 2. Lấy Time chính xác
	timestamp := time.Now().UnixNano()

	// 3. Tạo chuỗi ký
	rawString := fmt.Sprintf("%s|%s|%s|%d|%d|%s",
		signerID, docNum, action, stepOrder, timestamp, dataHash)

	// 4. HMAC Hash
	h := hmac.New(sha256.New, []byte(s.SecretKey.Secret))
	h.Write([]byte(rawString))
	signatureHash := hex.EncodeToString(h.Sum(nil))

	return signatureHash, dataHash, timestamp
}
