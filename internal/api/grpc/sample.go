// Initialize gRPC server
// internal/api/grpc/grpc.sample.go
package grpc

import "google.golang.org/grpc"

// NewServer returns a gRPC server ready for service registration.
// The svc parameter is accepted as interface{} to avoid coupling here;
// replace with the concrete service type if you prefer.
func NewServer(svc interface{}) *grpc.Server {
	s := grpc.NewServer()

	// Register your generated protobuf service implementations here, for example:
	// import pb "github.com/bharathbbg/inventory-service/internal/api/pb"
	// pb.RegisterInventoryServiceServer(s, NewInventoryGRPCHandler(svc.(*service.InventoryService)))

	// You can provide a handler adapter like:
	// type inventoryHandler struct {
	//     pb.UnimplementedInventoryServiceServer
	//     svc *service.InventoryService
	// }
	// func NewInventoryGRPCHandler(svc *service.InventoryService) *inventoryHandler { return &inventoryHandler{svc: svc} }
	// and implement the RPC methods on inventoryHandler.

	return s
}
