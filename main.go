package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "github.com/navdeep-singh-ghotra/kubecombe" // Your Protobuf definitionS
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{ID: "1", Name: "Laptop", Price: 999.99},
	{ID: "2", Name: "Phone", Price: 699.99},
}

// REST Handler
func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// gRPC Server
type server struct{ pb.UnimplementedProductServiceServer }

func (s *server) GetProducts(_ *pb.Empty, stream pb.ProductService_GetProductsServer) error {
	for _, p := range products {
		if err := stream.Send(&pb.Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: float32(p.Price),
		}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Start REST server
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/products", getProducts)
	go func() { log.Fatal(r.Run(":8080")) }()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{})
	log.Fatal(s.Serve(lis))
}