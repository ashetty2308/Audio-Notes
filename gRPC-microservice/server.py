import service_definition_pb2
import service_definition_pb2_grpc
import grpc
from concurrent import futures

class TextExtractionService(service_definition_pb2_grpc.TextExtractionServiceServicer):
    def ExtractText(self, request, context):
        s3_url = request.s3_url
        text = str(s3_url)
        return service_definition_pb2.ExtractTextResponse(extracted_text=text)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_definition_pb2_grpc.add_TextExtractionServiceServicer_to_server(TextExtractionService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()