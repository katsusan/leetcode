protoc  -ID:\Projects\protoc-3.6.1-win32\include --go_out=D:\Projects\protoc-3.6.1-win32\addressbook  D:\Projects\protoc-3.6.1-win32\addressbook.proto

# 1. 简易入门

example：
> server侧
    lis, err := net.Listen("tcp", port)
    s := grpc.NewServer()
    pb.RegisterGreeterServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        ...
    }

> client侧
    conn, err := grpc.Dial(address, grpc.WithInsecure())    //address为server侧的监听地址
    defer conn.Close()
    c := pb.NewGreeterClient(conn)  //*.pb.go为protoc的编译结果
    r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})

# 2. 前置知识protobuffer



