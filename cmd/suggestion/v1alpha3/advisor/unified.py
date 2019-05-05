import grpc
from concurrent import futures

import time

from pkg.api.suggestion.v1alpha3.python import advisor_pb2_grpc
from pkg.suggestion.v1alpha3.advisor.advisor import AdvisorUnifiedService
from pkg.suggestion.v1alpha1.types import DEFAULT_PORT

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    advisor_pb2_grpc.add_AdvisorSuggestionServicer_to_server(
        AdvisorUnifiedService(), server)
    server.add_insecure_port(DEFAULT_PORT)
    print("Listening...")
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()
