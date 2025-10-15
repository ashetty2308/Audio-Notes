import service_definition_pb2
import service_definition_pb2_grpc
import grpc
from concurrent import futures
import os
import boto3
import fitz
from io import BytesIO
from transformers import pipeline
from dotenv import load_dotenv
from pathlib import Path

summarizer = pipeline("summarization", model="facebook/bart-large-cnn", device=-1)

def get_bucket_key_filetype(url):
    prefix = "s3://"
    # postfix will represent bucket/filepath/as/in/bucket
    postfix = url[len(prefix):]
    # maxsplit = 1 here since we only want to split once into 2 buckets, key and file (could be nested file structure)
    bucket, key = postfix.split("/", 1)
    filetype = key.split('.', 1)[1]
    return bucket, key, filetype

def summarize_text(all_text):
    result = summarizer(all_text, max_length=100, min_length=20, do_sample=False)
    return result[0]["summary_text"]

class TextExtractionService(service_definition_pb2_grpc.TextExtractionServiceServicer):
    def ExtractText(self, request, context):
        s3_url = request.s3_url
        if not s3_url.startswith("s3://"):
            raise ValueError("Invalid url! (doesn't begin with correct s3:// path)")
        bucket, key, filetype = get_bucket_key_filetype(s3_url)
        s3 = boto3.client('s3')
        file_bytes = s3.get_object(Bucket=bucket, Key=key)['Body'].read()
        with fitz.open(stream=BytesIO(file_bytes), filetype=filetype) as notes:
            text = ""
            for page in notes:
                text += page.get_text()
        return service_definition_pb2.ExtractTextResponse(extracted_text=summarize_text(text))
      
def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_definition_pb2_grpc.add_TextExtractionServiceServicer_to_server(TextExtractionService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
