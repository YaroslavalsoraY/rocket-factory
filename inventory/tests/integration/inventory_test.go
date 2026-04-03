package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventory_v1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventory_v1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearPartsCollection(ctx)

		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("Get", func() {
		var uuid string
		BeforeEach(func() {
			var err error
			uuid, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых данных в MongoDB")
		})
		It("должен успешно возвращать информацию о детали по uuid", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventory_v1.GetPartRequest{
				Uuid: uuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart().Uuid).ToNot(BeEmpty())
		})
	})

	Describe("List", func() {
		uuids := make([]string, 0)
		BeforeEach(func() {
			for range testPartsCount {
				uuid, err := env.InsertTestPart(ctx)
				uuids = append(uuids, uuid)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых данных в MongoDB")
			}
		})
		It("должен возвращать список всех деталей", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{Filter: &inventory_v1.PartsFilter{}})

			Expect(err).ToNot(HaveOccurred())

			result := resp.GetParts()

			Expect(len(result)).To(Equal(len(uuids)))
		})
	})
})
