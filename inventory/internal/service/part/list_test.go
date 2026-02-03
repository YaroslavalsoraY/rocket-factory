package part

import (
	"github.com/brianvoe/gofakeit/v6"
	"inventory/internal/model"
)

func (s *ServiceSuite) TestList() {
	var (
		testFilters = model.Filters{}
		uuid        = gofakeit.UUID()
		testParts   = []model.PartInfo{
			{
				UUID:          uuid,
				Name:          gofakeit.MinecraftArmorPart(),
				Description:   gofakeit.Paragraph(3, 5, 5, " "),
				Price:         gofakeit.Float64(),
				StockQuantity: gofakeit.Int64(),
				Category:      model.CategoryEnum_CATEGORY_PORTHOLE,
				Tags:          []string{"тормоза", "диск", "передний", "вентилируемый"},
			},
		}
	)

	s.InvRepository.On("ListParts", s.ctx, testFilters).Return(testParts, nil)

	res, err := s.service.List(s.ctx, model.Filters{})

	s.NoError(err)
	s.Equal(testParts, res)
}

func (s *ServiceSuite) TestListError() {
	testFilters := model.Filters{
		UUIDs: []string{gofakeit.UUID()},
	}

	s.InvRepository.On("ListParts", s.ctx, testFilters).Return([]model.PartInfo{}, model.ErrPartsNotFound)

	_, err := s.service.List(s.ctx, testFilters)

	s.ErrorIs(err, model.ErrPartsNotFound)
}
