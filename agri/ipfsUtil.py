from config import *
import requests
import shutil

def ipfsAdd(fileName):
    files = {'media': open(fileName, 'rb')}
    addUrl = "https://ipfs.infura.io:5001/api/v0/add?pin=false"
    headers = {
        "Content-Type: multipart/form-data"
    }
    print("request")
    r = requests.post(addUrl,files=files)
    print("request finish")
    return r.text
def ipfsCat(fileHash,outputFile = False):
    catUrl = "https://ipfs.infura.io:5001/api/v0/cat?arg=" + fileHash
    r = requests.get(catUrl, stream=True)
    if r.status_code != 200:
        print("ERROR, it may take some time or there isn't any node keep this file")
    if outputFile:
        with open("./outFile.jpg", 'wb') as f:
            f.write(r.content)
    return r.content
# ipfsAdd("test1.jpg")
# ipfsCat("QmVPEapUax45br7xqKkSTpxQixkgRpr1aMSDrvycL8QNA6")
