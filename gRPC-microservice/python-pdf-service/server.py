import service_definition_pb2
import service_definition_pb2_grpc
import grpc
from concurrent import futures
import os
import boto3


def get_bucket_and_file_name(url):
        prefix = "s3://"
        # postfix will represent bucket/filepath/as/in/bucket
        postfix = url[len(prefix):]
        # maxsplit = 1 here since we only want to split once into 2 buckets, key and file (could be nested file structure)
        bucket, key = postfix.split("/", 1)
        return bucket, key

class TextExtractionService(service_definition_pb2_grpc.TextExtractionServiceServicer):
            
    def ExtractText(self, request, context):
        s3_url = request.s3_url
        if not s3_url.startswith("s3://"):
            raise ValueError("Invalid url! (doesn't begin with correct s3:// path)")
        bucket, key = get_bucket_and_file_name(s3_url)
        s3 = boto3.client('s3')
        concat = bucket + ", " + key
        return service_definition_pb2.ExtractTextResponse(extracted_text=concat)
      
def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_definition_pb2_grpc.add_TextExtractionServiceServicer_to_server(TextExtractionService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()