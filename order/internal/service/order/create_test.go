package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestCreateOrder() {
	var (
		uuid      = gofakeit.UUID()
		testUUIDs = []string{uuid}
		userUUID  = gofakeit.UUID()
		testPrice = gofakeit.Float32Range(1, 99999)
		testParts = []*model.PartInfo{
			{
				UUID:          uuid,
				Name:          gofakeit.MinecraftArmorPart(),
				Description:   gofakeit.Paragraph(3, 5, 5, " "),
				Price:         float64(testPrice),
				StockQuantity: gofakeit.Int64(),
				Category:      model.CategoryEnum_CATEGORY_PORTHOLE,
				Tags:          []string{"тормоза", "диск", "передний", "вентилируемый"},
			},
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, model.Filters{UUIDs: testUUIDs}).Return(testParts, nil)

	s.orderRepository.On("CreateOrder", s.ctx, mock.Anything).Return(nil)

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, testUUIDs)
	s.NoError(err)
	s.Equal(testPrice, totalPrice)
	s.NotEmpty(orderUUID)
}

func (s *ServiceSuite) TestCreateOrderPartsNotFound() {
	var (
		uuid      = gofakeit.UUID()
		testUUIDs = []string{uuid}
		userUUID  = gofakeit.UUID()
		testParts = []*model.PartInfo{}
	)

	s.inventoryClient.On("ListParts", s.ctx, model.Filters{UUIDs: testUUIDs}).Return(testParts, nil)

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, testUUIDs)
	s.ErrorIs(err, model.ErrPartsNotFound)
	s.Equal(float32(0), totalPrice)
	s.Empty(orderUUID)
}
