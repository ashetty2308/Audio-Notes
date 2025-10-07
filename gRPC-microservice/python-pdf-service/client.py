import grpc
import service_definition_pb2
import service_definition_pb2_grpc

channel = grpc.insecure_channel('localhost:50051')
stub = service_definition_pb2_grpc.TextExtractionServiceStub(channel)  # must match proto service name

request = service_definition_pb2.ExtractTextRequest(s3_url="s3://s3-golang-uploaded-pdfs/Biology Notes.pdf")
response = stub.ExtractText(request)

print(response.extracted_text)
