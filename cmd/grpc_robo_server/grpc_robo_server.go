package main

import (
	"log"
	"net"

	pb "github.com/ClickerAI/ClickerAI/proto"

	"google.golang.org/grpc"

	"context"

	"github.com/go-vgo/robotgo"
)

func int32SliceToIntSlice(int32Slice []int32) []int {
	intSlice := make([]int, len(int32Slice))
	for i, v := range int32Slice {
		intSlice[i] = int(v)
	}
	return intSlice
}

type server struct {
	pb.UnimplementedRobotGoServer
}

func (s *server) MoveSmooth(ctx context.Context, req *pb.MoveSmoothRequest) (*pb.MoveSmoothResponse, error) {
	x, y := req.GetX(), req.GetY()
	args := req.GetArgs()
	success := robotgo.MoveSmooth(int(x), int(y), args)

	return &pb.MoveSmoothResponse{Success: success}, nil
}

func (s *server) MoveRelative(ctx context.Context, req *pb.MoveRelativeRequest) (*pb.MoveRelativeResponse, error) {
	x, y := req.GetX(), req.GetY()
	robotgo.MoveRelative(int(x), int(y))
	return &pb.MoveRelativeResponse{}, nil
}

func (s *server) MoveSmoothRelative(ctx context.Context, req *pb.MoveSmoothRelativeRequest) (*pb.MoveSmoothRelativeResponse, error) {
	x, y := req.GetX(), req.GetY()
	args := req.GetArgs()
	robotgo.MoveSmoothRelative(int(x), int(y), args)

	return &pb.MoveSmoothRelativeResponse{}, nil
}

func (s *server) GetMousePos(ctx context.Context, req *pb.GetMousePosRequest) (*pb.GetMousePosResponse, error) {
	x, y := robotgo.GetMousePos()
	return &pb.GetMousePosResponse{X: int32(x), Y: int32(y)}, nil
}

func (s *server) Click(ctx context.Context, req *pb.ClickRequest) (*pb.ClickResponse, error) {
	args := req.GetArgs()
	robotgo.Click(args)
	return &pb.ClickResponse{}, nil
}

func (s *server) Scroll(ctx context.Context, req *pb.ScrollRequest) (*pb.ScrollResponse, error) {
	x, y := req.GetX(), req.GetY()
	args := req.GetArgs()

	robotgo.Scroll(int(x), int(y), int32SliceToIntSlice(args)...)

	return &pb.ScrollResponse{}, nil
}

func (s *server) ScrollMouse(ctx context.Context, req *pb.ScrollMouseRequest) (*pb.ScrollMouseResponse, error) {
	x := req.GetX()
	direction := req.GetDirection()
	robotgo.ScrollMouse(int(x), direction...)

	return &pb.ScrollMouseResponse{}, nil
}

func (s *server) ScrollSmooth(ctx context.Context, req *pb.ScrollSmoothRequest) (*pb.ScrollSmoothResponse, error) {
	to := req.GetTo()
	args := req.GetArgs()
	robotgo.ScrollSmooth(int(to), int32SliceToIntSlice(args)...)
	return &pb.ScrollSmoothResponse{}, nil
}

func (s *server) ScrollRelative(ctx context.Context, req *pb.ScrollRelativeRequest) (*pb.ScrollRelativeResponse, error) {
	x, y := req.GetX(), req.GetY()
	args := req.GetArgs()

	robotgo.ScrollRelative(int(x), int(y), int32SliceToIntSlice(args)...)
	return &pb.ScrollRelativeResponse{}, nil
}

func (s *server) MilliSleep(ctx context.Context, req *pb.MilliSleepRequest) (*pb.MilliSleepResponse, error) {
	tm := req.GetTm()
	robotgo.MilliSleep(int(tm))
	return &pb.MilliSleepResponse{}, nil
}

func (s *server) Sleep(ctx context.Context, req *pb.SleepRequest) (*pb.SleepResponse, error) {
	tm := req.GetTm()
	robotgo.Sleep(int(tm))
	return &pb.SleepResponse{}, nil
}

func (s *server) CaptureScreen(ctx context.Context, req *pb.CaptureScreenRequest) (*pb.CaptureScreenResponse, error) {
	args := req.GetArgs()
	// bitmap :=
	robotgo.CaptureScreen(int32SliceToIntSlice(args)...)

	return &pb.CaptureScreenResponse{
		// Bitmap: bitmap
	}, nil
}

func (s *server) CaptureGo(ctx context.Context, req *pb.CaptureGoRequest) (*pb.CaptureGoResponse, error) {
	args := req.GetArgs()
	// bitmap :=
	robotgo.CaptureGo(int32SliceToIntSlice(args)...)

	return &pb.CaptureGoResponse{
		// Bitmap: bitmap
	}, nil

}

func (s *server) CaptureImg(ctx context.Context, req *pb.CaptureImgRequest) (*pb.CaptureImgResponse, error) {
	args := req.GetArgs()
	// bitmap :=
	robotgo.CaptureImg(int32SliceToIntSlice(args)...)

	return &pb.CaptureImgResponse{
		// Bitmap: bitmap
	}, nil
}

func (s *server) Move(ctx context.Context, req *pb.MoveRequest) (*pb.MoveResponse, error) {
	x, y := req.GetX(), req.GetY()
	robotgo.Move(int(x), int(y))
	return &pb.MoveResponse{}, nil
}

func (s *server) DragSmooth(ctx context.Context, req *pb.DragSmoothRequest) (*pb.DragSmoothResponse, error) {
	x, y := req.GetX(), req.GetY()
	args := req.GetArgs()
	robotgo.DragSmooth(int(x), int(y), args)
	return &pb.DragSmoothResponse{}, nil
}

func main() {
	// create listener
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	pb.RegisterRobotGoServer(s, &server{})

	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
