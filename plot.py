import numpy as np
import matplotlib as mpl
import matplotlib.pyplot as plt
import json

podname = "dfec999007fe7011eb0ad0e08110ed4eba5f-yangzh321-0"

def plotcpu():
    file = open("./result/cpuinfo.json", 'r', encoding='utf-8')
    for line in file.readlines():
        dic = json.loads(line)
        if dic['pod'] != podname:
            continue
        val = dic['history']
        if val == None:
            val = []
        plt.title(dic['pod'])
        
        plt.xlim(0,len(val) if len(val)!=0 else 1)
        plt.xlabel('Time/min')
        plt.ylabel('Cores')
        plt.plot(val)
        plt.tight_layout()
        plt.savefig("./cpu-"+dic['pod']+'.jpg')
        #plt.show()
        plt.close()
    file.close()

def plotpodgpu():
    file = open("./result/gpuinfoutil.json", 'r', encoding='utf-8')
    poduuiduti = {}
    for line in file.readlines():
        dic = json.loads(line)
        val = dic['history']
        if val == None:
            val = []
        uuid = poduuiduti.get(dic['pod'],{})
        uuid[dic['uuid']]=val
        poduuiduti[dic['pod']] = uuid
    file.close()

    file = open("./result/gpuinfomem.json", 'r', encoding='utf-8')
    poduuidmem = {}
    for line in file.readlines():
        dic = json.loads(line)
        val = dic['history']
        if val == None:
            val = []
        uuid = poduuidmem.get(dic['pod'],{})
        uuid[dic['uuid']]=val
        poduuidmem[dic['pod']] = uuid

    for pod in poduuidmem:
        gpucnt = len(poduuidmem[pod])
        cnt = 0
        col = 1
        row = 1
        if gpucnt == 1:
            row = 1
        elif gpucnt <= 2 :
            row = 2
        elif gpucnt <= 4 :
            col = 2
            row = 2
        elif gpucnt <= 8 :
            col = 4
            row = 2
        elif gpucnt >= 16:
            row = 4
            col = gpucnt / 4
        if pod != podname:
            continue
        _, axes = plt.subplots(row,col,squeeze=False)
        plt.title(pod)
        print(pod)
        for uuid in poduuidmem[pod]:
            idxx=int(cnt/col)
            idxy=int(cnt%col)
            ax = axes[idxx][idxy]
            ax.set_title(uuid)
            ax.set_xlabel('Time/min')
            ax.set_ylabel('GPU Mem/MB')
            #ax.set_ylim(bottom=0)
            val1 = poduuidmem.get(pod,{}).get(uuid,{})
            length = len(val1) if len(val1) !=0 else 1
            x = np.arange(0,float(length)*0.5,0.5)
            lns1=ax.plot(x[:300],val1[:300],c='r',label="GPU Mem")
            #lns1=ax.plot(x,val1,c='r',label="GPU Mem")

            ax2=ax.twinx()
            ax2.set_ylabel('GPU Util/%')
            #ax2.set_ylim(bottom=0)
            val2 = poduuiduti.get(pod,{}).get(uuid,{})

            lns2=ax2.plot(x[:300],val2[:300],c='g')
            #lns2=ax2.plot(x,val2,c='g')
            
            lns=lns1+lns2
            plt.legend(lns,["GPU Mem","GPU Util"],loc=0)
            cnt = cnt+1
        plt.suptitle(pod)
        
        plt.tight_layout()
        plt.savefig("./gpu-"+pod+'.jpg')
        #plt.show()
        plt.close()
    file.close()

def plotmem():
    file = open("./result/meminfo.json", 'r', encoding='utf-8')
    for line in file.readlines():
        dic = json.loads(line)
        if dic['pod'] != podname:
            continue
        val = dic['history']
        if val == None:
            val = []
        
        length = len(val) if len(val) !=0 else 1
        x = np.arange(0,float(length)*0.5,0.5)
        # plt.xlim(0,length)
        plt.title(dic['pod'])
        plt.xlabel('Time/min')
        plt.ylabel('Mem/B')
        plt.plot(x,val)
        plt.tight_layout()
        plt.savefig("./mem-"+dic['pod']+'.jpg')
        # plt.show()
        # break
        plt.close()
    file.close()

def plotnodecpu():
    file = open("./result/nodecpuuti.json", 'r', encoding='utf-8')
    for line in file.readlines():
        dic = json.loads(line)
        if dic['pod'] != podname:
            continue
        val = dic['utili']
        if val == None:
            val = []
        plt.title(dic['pod'])
        
        length = len(val) if len(val) !=0 else 1
        x = np.arange(0,float(length)*0.5,0.5)
        # plt.xlim(0,length)
        plt.xlabel('Time/min')
        plt.ylabel('CPU Util/%')
        plt.plot(x,val)
        plt.tight_layout()
        plt.savefig("./cpu-"+dic['pod']+'.jpg')
        # plt.show()
        # break
        plt.close()
    file.close()

def plotnodegpu():
    file = open("./result/nodegpuinfo.json", 'r', encoding='utf-8')
    for line in file.readlines():
        dic = json.loads(line)
        state = dic['state']
        val = []
        for v in state:
            val.append(v['use']*100.0/v['total'])

        plt.xlim(0,len(val) if len(val)!=0 else 1)
        plt.title(dic['name'])
        plt.xlabel('Time/min')
        plt.ylabel('Node GPU Util/%')
        plt.plot(val)
        plt.tight_layout()
    #    plt.savefig("./plot/nodegpuinfo/"+dic['name']+'.jpg')
        plt.show()
        plt.close()
    file.close()

def plotIO():
    file = open("./result/nodereiverate.json", 'r', encoding='utf-8')
    nodeio = {}
    for line in file.readlines():
        dic = json.loads(line)
        node = dic['node']
        val = dic['Irate']
        if val == None:
                val = []
        nodeio[dic['node']] = val
        
        # plt.xlim(0,len(val) if len(val)!=0 else 1)
        # plt.title(dic['node'])
        # plt.xlabel('Time/min')
        # plt.ylabel('Receive Byte Rate/B/s')
        # plt.plot(val)
        # #plt.savefig("./plot/nodereiverate/"+dic['node']+'.jpg')
        # plt.show()
    file.close()

    file = open("./result/nodetransrate.json", 'r', encoding='utf-8')
    for line in file.readlines():
        dic = json.loads(line)
        node = dic['node']
        val = dic['Orate']
        # if node != "agx-16":
        #     continue
        _, ax = plt.subplots()
        
        ax.set_title(node)
        ax.set_xlabel('Time/min')
        ax.set_ylabel('Recieve rate/B/s')
        val1 = nodeio.get(node,{})
        length = len(val1) if len(val1) !=0 else 1
        x = np.arange(length)
        lns1=ax.plot(x,val1,c='r')
            
        ax2=ax.twinx()
        ax2.set_ylabel('Transmit Rate/B/s')
        lns2=ax2.plot(x,val,c='g')
        lns=lns1+lns2
        plt.legend(lns,["Recieve rate","Transmit Rate"],loc=0)
        plt.tight_layout()
        plt.savefig("./plot/noderate/"+dic['node']+'.jpg')
        #plt.show()
        plt.close()
    file.close()

def plotIBIO():
    file = open("./result/nodeibreiverate.json", 'r', encoding='utf-8')
    nodeio = {}
    for line in file.readlines():
        dic = json.loads(line)
        node = dic['node']
        val = dic['Irate']
        if val == None:
            val = []
        nodeio[dic['node']] = val
        
        # plt.xlim(0,len(val) if len(val)!=0 else 1)
        # plt.title(dic['node'])
        # plt.xlabel('Time/min')
        # plt.ylabel('Receive Byte Rate/B/s')
        # plt.plot(val)
        # #plt.savefig("./plot/nodereiverate/"+dic['node']+'.jpg')
        # plt.show()
    file.close()

    file = open("./result/nodeibtransrate.json", 'r', encoding='utf-8')
    for line in file.readlines():
        dic = json.loads(line)
        node = dic['node']
        val = dic['Orate']
        if val == None:
            val = []
        _, ax = plt.subplots()
        
        ax.set_title(node)
        ax.set_xlabel('Time/min')
        ax.set_ylabel('IB Recieve rate/B/s')
        val1 = nodeio.get(node,{})
        length = len(val1)
        x = np.arange(length)
        lns1=ax.plot(x,val1,c='r')
            
        ax2=ax.twinx()
        ax2.set_ylabel('IB Transmit Rate/B/s')
        lns2=ax2.plot(x,val,c='g')
        lns=lns1+lns2
        plt.legend(lns,["IB Recieve rate","IB Transmit Rate"],loc=0)
        plt.tight_layout()
        plt.savefig("./plot/nodeib/"+dic['node']+'.jpg')
        #plt.show()
        plt.close()
    file.close()    

def cdf():
    file = open("./2.log", 'r', encoding='utf-8')
    str = file.read()
    data = str.split(" ")
    tot = []
    cnt = []
    mus = 0.0
    tasknum = 752
    for v in data :
        mus = mus + float(v)
        cnt.append(int(v))
        tot.append(mus*100/tasknum)
        
    fig, ax = plt.subplots()
    plt.title("GPU Uti Lesser 10% CPU CDF")
    ax.set_xlabel('CPU Utili')
    ax.set_ylabel('Task Count')
    # ax.set_ylim(bottom=0)
    # cnt[100]=150
    # cnt[1]=110

    xx = np.arange(0,101)
    ax.bar(xx,cnt)
    # plt.text(xx[100],cnt[100],372,ha='center')
    # plt.xticks(rotation=45)
    # plt.text(xx[1],cnt[1],299,ha='center')
    # plt.xticks(rotation=45)
    ax2=ax.twinx()
    ax2.set_ylabel('Occupy/%')
    #ax2.set_ylim(bottom=0)
    ax2.plot(xx,tot,c="orange")

    plt.tight_layout()
    plt.savefig("./2.jpg")
    #plt.show()
    #print(tot)
    file.close()
    plt.close()

def main():
    # plotcpu()
    plotpodgpu()
    plotmem()
    plotnodecpu()
    # plotnodegpu()
    # plotIO()
    # plotIBIO()
    # cdf()

if __name__ == '__main__':
    main()