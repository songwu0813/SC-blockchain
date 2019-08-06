from config import *
import subprocess
import json
import base58
import os
import requests
from web3 import Web3

w3 = Web3(Web3.HTTPProvider(ethereumNode))
if not w3.isConnected():
    print("ethereum node error")
contract_ = w3.eth.contract(
abi=contract_interface['abi'],
bytecode=contract_interface['bin'])
acct = w3.eth.account.privateKeyToAccount(privateKey)
construct_txn = contract_.constructor().buildTransaction({
'from': acct.address,
'nonce': w3.eth.getTransactionCount(acct.address),
'gas': 1728712,
'gasPrice': w3.toWei('21', 'gwei')})
signed = acct.signTransaction(construct_txn)
temp = w3.eth.sendRawTransaction(signed.rawTransaction)
print("txid : ", w3.toHex(temp))
