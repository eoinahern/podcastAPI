package mocks

import "github.com/eoinahern/podcastAPI/models"

type MockMailRequest struct {
	SenderId     string
	toId         string
	BodyLocation string
	bodyParams   *models.TemplateParams
}

//SendMail send fake mail
func (m *MockMailRequest) SendMail() (bool, error) {
	return true, nil
}

//SetBodyParams mockBodyParams
func (m *MockMailRequest) SetBodyParams(bodyParams *models.TemplateParams) {
	m.bodyParams = bodyParams
}

//SetToID set toId
func (m *MockMailRequest) SetToID(toid string) {
	m.toId = toid
}

// MockPasswordEncryptUitl mock implementation
type MockPasswordEncryptUitl struct {
}

func (m *MockPasswordEncryptUitl) Encrypt(password string) string {
	return password
}

func (m *MockPasswordEncryptUitl) CheckSame(dbPass string, sentPass string) bool {
	return true
}

type MockFileHelperUtil struct {
}

func (m *MockFileHelperUtil) CheckDirFileExists(path string) bool {

	if path == "../debug_files/1/podcast" {
		return true
	}

	return false
}

func (m *MockFileHelperUtil) CreateDir(dirpath string) {

}
