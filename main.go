package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "path/to/proto"  // Your Protobuf definition
)

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	products = []Product{{ID: "1", Name: "Laptop"}}
	mu       sync.Mutex
)

// REST Handler
func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// gRPC Handler
type server struct{ pb.UnimplementedProductServiceServer }
func (s *server) GetProducts(req *pb.Empty, stream pb.ProductService_GetProductsServer) error {
	for _, p := range products {
		if err := stream.Send(&pb.Product{Id: p.ID, Name: p.Name}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// REST Server
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/products", getProducts)
	go func() { log.Fatal(r.Run(":8080")) }()

	// gRPC Server
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{})
	log.Fatal(s.Serve(lis))
}