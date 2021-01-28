# sample suggestion server for darts, 
# which should written by user and run in suggestion resource

import grpc
import time
from concurrent import futures

from pkg.api.suggestion.v1alpha3.python import oneshot_pb2_grpc
from pkg.api.suggestion.v1alpha3.python import oneshot_pb2
from pkg.suggestion.v1alpha1.types import DEFAULT_PORT
from cnn import train

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
CHUNK_SIZE = 64 * 1024 # 64 K
LISTEN_ADDRESS = '127.0.0.1:6789'


class DartsService(oneshot_pb2_grpc.OneshotSuggestionServicer):
    def __init__(self):
        self.train_ = train.DartsTrain()

    def GetSuggestions(self, request, context):
        # do the training (fake) and export darts            
        self.train_.export_darts()

        with open("./model/darts.proto", "rb") as f:
            content = f.read(CHUNK_SIZE)
            yield oneshot_pb2.GetOneshotSuggestionReply(onnx_model=content)

            while content != b"":
                # Do stuff with byte.
                content = f.read(CHUNK_SIZE)
                yield oneshot_pb2.GetOneshotSuggestionReply(onnx_model=content)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    oneshot_pb2_grpc.add_OneshotSuggestionServicer_to_server(DartsService(), server)
    server.add_insecure_port(LISTEN_ADDRESS)
    print("Listening...")
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()


        
      
