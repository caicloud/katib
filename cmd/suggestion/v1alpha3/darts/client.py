# sample client for GRPC getSuggestions
import grpc
import time
from concurrent import futures

from pkg.api.suggestion.v1alpha3.python import oneshot_pb2_grpc
from pkg.api.suggestion.v1alpha3.python import oneshot_pb2

ADDRESS = '127.0.0.1:6789'

def main():
    channel = grpc.insecure_channel(ADDRESS)
    stub = oneshot_pb2_grpc.OneshotSuggestionStub(channel)
    request = oneshot_pb2.GetOneshotSuggestionRequest()
    reply = stub.GetSuggestions(request)
    cnt = 0
    for reply_item in reply:
        cnt += 1
        print(reply_item)
    print(cnt) # how many chunks are received

if __name__ == "__main__":
    main()
