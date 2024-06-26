# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import pb2s.profile_pb2 as profile__pb2


class ProfileServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CreateProfile = channel.unary_unary(
                '/proto.ProfileService/CreateProfile',
                request_serializer=profile__pb2.ProfileRequest.SerializeToString,
                response_deserializer=profile__pb2.ErrorResponse.FromString,
                )
        self.ReadProfile = channel.unary_unary(
                '/proto.ProfileService/ReadProfile',
                request_serializer=profile__pb2.IdRequest.SerializeToString,
                response_deserializer=profile__pb2.Profile.FromString,
                )
        self.UpdateProfile = channel.unary_unary(
                '/proto.ProfileService/UpdateProfile',
                request_serializer=profile__pb2.ProfileRequest.SerializeToString,
                response_deserializer=profile__pb2.ErrorResponse.FromString,
                )
        self.GetNextProfile = channel.unary_unary(
                '/proto.ProfileService/GetNextProfile',
                request_serializer=profile__pb2.IdRequest.SerializeToString,
                response_deserializer=profile__pb2.Profile.FromString,
                )
        self.Like = channel.unary_unary(
                '/proto.ProfileService/Like',
                request_serializer=profile__pb2.TargetRequest.SerializeToString,
                response_deserializer=profile__pb2.ErrorResponse.FromString,
                )
        self.GetLikes = channel.unary_unary(
                '/proto.ProfileService/GetLikes',
                request_serializer=profile__pb2.Empty.SerializeToString,
                response_deserializer=profile__pb2.LikesResponse.FromString,
                )


class ProfileServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def CreateProfile(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ReadProfile(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UpdateProfile(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetNextProfile(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Like(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetLikes(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ProfileServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CreateProfile': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateProfile,
                    request_deserializer=profile__pb2.ProfileRequest.FromString,
                    response_serializer=profile__pb2.ErrorResponse.SerializeToString,
            ),
            'ReadProfile': grpc.unary_unary_rpc_method_handler(
                    servicer.ReadProfile,
                    request_deserializer=profile__pb2.IdRequest.FromString,
                    response_serializer=profile__pb2.Profile.SerializeToString,
            ),
            'UpdateProfile': grpc.unary_unary_rpc_method_handler(
                    servicer.UpdateProfile,
                    request_deserializer=profile__pb2.ProfileRequest.FromString,
                    response_serializer=profile__pb2.ErrorResponse.SerializeToString,
            ),
            'GetNextProfile': grpc.unary_unary_rpc_method_handler(
                    servicer.GetNextProfile,
                    request_deserializer=profile__pb2.IdRequest.FromString,
                    response_serializer=profile__pb2.Profile.SerializeToString,
            ),
            'Like': grpc.unary_unary_rpc_method_handler(
                    servicer.Like,
                    request_deserializer=profile__pb2.TargetRequest.FromString,
                    response_serializer=profile__pb2.ErrorResponse.SerializeToString,
            ),
            'GetLikes': grpc.unary_unary_rpc_method_handler(
                    servicer.GetLikes,
                    request_deserializer=profile__pb2.Empty.FromString,
                    response_serializer=profile__pb2.LikesResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'proto.ProfileService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ProfileService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def CreateProfile(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.ProfileService/CreateProfile',
            profile__pb2.ProfileRequest.SerializeToString,
            profile__pb2.ErrorResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ReadProfile(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.ProfileService/ReadProfile',
            profile__pb2.IdRequest.SerializeToString,
            profile__pb2.Profile.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def UpdateProfile(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.ProfileService/UpdateProfile',
            profile__pb2.ProfileRequest.SerializeToString,
            profile__pb2.ErrorResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetNextProfile(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.ProfileService/GetNextProfile',
            profile__pb2.IdRequest.SerializeToString,
            profile__pb2.Profile.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Like(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.ProfileService/Like',
            profile__pb2.TargetRequest.SerializeToString,
            profile__pb2.ErrorResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetLikes(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/proto.ProfileService/GetLikes',
            profile__pb2.Empty.SerializeToString,
            profile__pb2.LikesResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
