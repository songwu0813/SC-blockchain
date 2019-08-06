from web3 import Web3
import sys
import requests
import os
import base58
import json
import subprocess
from pprint import pprint
from config import *
import hashlib
from ipfsUtil import ipfsAdd, ipfsCat
w3 = Web3(Web3.HTTPProvider(ethereumNode))
if not w3.isConnected():
    print("ethereum node error")
    exit()


def convertBytes32Ipfs(b):
    prefix = b'\x12 '
    b = prefix + b
    b = base58.b58encode(b)
    return b


def convertIpfsBytes32(hash_string):
    b = base58.b58decode(hash_string)
    ## print("base58 `hex` of ",hash_string,"=",b.hex())
    return b[2:]


if len(sys.argv) == 3:
    if sys.argv[1].lower() == "add":
        # ipfs start
        fileInfo = (sys.argv[2]).split(":")
        id = fileInfo[0]
        ## print("ready to add ", id,"'s file:",fileInfo[1]," to ipfs")
        ## fileHash = ipfsAdd(fileInfo[1])
        ## print(fileHash)
        ## picHash = json.loads(fileHash)['Hash']
        ## picHash = convertIpfsBytes32(picHash)
        ## print(picHash.hex())
        filename = fileInfo[1]
        with open(filename,"rb") as f:
            bytes = f.read() # read entire file as bytes
            readable_hash = hashlib.sha256(bytes).hexdigest();
            print(readable_hash)

        # etherum start
        ## print("ready to add fileHash to ethereum")
        contract_ = w3.eth.contract(
            address="0xecab3320Ca2d9377850428fe28C302573B8f0C16", abi=contract_interface["abi"])
        acct = w3.eth.account.privateKeyToAccount(privateKey)
        construct_txn = contract_.functions.addNewHash(id, readable_hash).buildTransaction({
            'from': acct.address,
            'nonce': w3.eth.getTransactionCount(acct.address),
            'gas': 1728712,
            'gasPrice': w3.toWei('21', 'gwei')})
        signed = acct.signTransaction(construct_txn)
        txId = w3.eth.sendRawTransaction(signed.rawTransaction)
        print(w3.toHex(txId))
    elif sys.argv[1].lower() == "cat":
        print("ready to get file from ipfs")
        ipfsCat(sys.argv[2], True)
        print("file output to ./outFile.jpg")
else:
    print("read reademe, and try it again")
