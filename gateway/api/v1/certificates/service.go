package certificates

type CertificateService struct {
	repo *CertificateRepository
}

func NewCertificateService(repo *CertificateRepository) *CertificateService {
	return &CertificateService{repo: repo}
}

func (s *CertificateService) GetMyCertificates(userID uint) ([]CertificateDTO, error) {
	certificares, err := s.repo.FindCertificatesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var dtos []CertificateDTO
	for _, cert := range certificares {
		courseName := ""
		if cert.Course != nil {
			courseName = cert.Course.Name
		}
		dto := CertificateDTO{
			CourseName: courseName,
			Code:       cert.Code,
			IssuedAt:   cert.IssuedAt,
		}
		dtos = append(dtos, dto)
	}

	if dtos == nil {
		dtos = []CertificateDTO{}
	}
	return dtos, nil
}
