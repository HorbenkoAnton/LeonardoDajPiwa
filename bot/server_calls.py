import pb2s.profile_pb2 as profile_pb2
import pb2s.profile_pb2_grpc as profile_pb2_grpc
import grpc


#Gets array profile
#Returns request report
def CreateRequest(profile):
    with grpc.insecure_channel("localhost:50051") as chan:
        stub = profile_pb2_grpc.ProfileServiceStub(chan)
        create_request = profile_pb2.ProfileRequest(profile=profile)
        return stub.CreateProfile(create_request)

def ReadRequest(id):
    with grpc.insecure_channel("localhost:50051") as chan:
            stub = profile_pb2_grpc.ProfileServiceStub(chan)
            read_request = profile_pb2.IdRequest(id=id)
            return stub.ReadProfile(read_request)
    
def UpdateRequest(profile):
    with grpc.insecure_channel("localhost:50051") as chan:
            stub = profile_pb2_grpc.ProfileServiceStub(chan)
            # profile = profile_pb2.Profile(
            #                         ID=profile_arg[0], 
            #                         name=profile_arg[1],
            #                         age=profile_arg[2],
            #                         description=profile_arg[3],
            #                         location =profile_arg[4],
            #                         )
        
            request = profile_pb2.ProfileRequest(profile=profile)
            return stub.UpdateProfile(request)
    
def GetNextProfileRequest(id):
    with grpc.insecure_channel("localhost:50051") as chan:
            stub = profile_pb2_grpc.ProfileServiceStub(chan)
            request = profile_pb2.IdRequest(id=id)
            return stub.GetNextProfile(request)
    

def LikeRequest(id, tgtId):
    with grpc.insecure_channel("localhost:50051") as chan:
        stub = profile_pb2_grpc.ProfileServiceStub(chan)
        request = profile_pb2.TargetRequest(id = id,tgtId=tgtId)
        return stub.Like(request)
    
def GetLikesRequest():
    with grpc.insecure_channel("localhost:50051") as chan:
        stub = profile_pb2_grpc.ProfileServiceStub(chan)
        request = profile_pb2.Empty()
        return stub.GetLikes(request)