package part

import (
	"inventory/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestGetPart() {
	var (
		uuid = gofakeit.UUID()
		testPart = model.PartInfo{
			UUID:          uuid,
			Name:          gofakeit.MinecraftArmorPart(),
			Description:   gofakeit.Paragraph(3, 5, 5, " "),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      model.CategoryEnum_CATEGORY_PORTHOLE,
			Tags:          []string{"тормоза", "диск", "передний", "вентилируемый"},
		}
	)

	s.InvRepository.On("GetPart", s.ctx, uuid).Return(testPart, nil)

	part, err := s.service.Get(s.ctx, uuid)

	s.NoError(err)
	s.Equal(testPart, part)
}

func (s *ServiceSuite) TestGetError() {
	var (
		uuid = gofakeit.UUID()
		emptyPart = model.PartInfo{}
	)

	s.InvRepository.On("GetPart", s.ctx, uuid).Return(emptyPart, model.ErrPartsNotFound)

	part, err := s.service.Get(s.ctx, uuid)

	s.ErrorIs(err, model.ErrPartsNotFound)
	s.Equal(emptyPart, part)
}

