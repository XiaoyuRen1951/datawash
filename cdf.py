import numpy as np
import matplotlib as mpl
import matplotlib.pyplot as plt
import json

def SCGPUUtui():
    file = open("./1.log", 'r', encoding='utf-8')
    str = file.read()
    data = str.split(" ")
    tot = []
    cnt = []
    mus = 0.0
    tasknum = 2612
    for v in data :
        mus = mus + float(v)
        cnt.append(int(v))
        tot.append(mus*100/tasknum)
    _, ax = plt.subplots()
    plt.title("Greater 3h Single Card TASK CDF ")
    ax.set_xlabel('GPU Util/%')
    ax.set_ylabel('Task Count')
    #ax.set_ylim(bottom=0)
    # cnt[0]=100

    xx = np.arange(len(data))
    ax.bar(xx,cnt)
    # plt.text(xx[0],cnt[0],506,ha='center')
    # plt.xticks(rotation=45)
    ax2=ax.twinx()
    ax2.set_ylabel('Occupy/%')
    #ax2.set_ylim(bottom=0)
    ax2.plot(xx,tot,c="orange")

    plt.tight_layout()
    plt.savefig("./plot/cdf-ALL.jpg")
    plt.show()
    #print(tot)
    file.close()

def MCGPUUti():
    file = open("./1.log", 'r', encoding='utf-8')
    str = file.read()
    data = str.split(" ")
    tot = []
    cnt = []
    mus = 0.0
    tasknum = 863
    for v in data :
        mus = mus + float(v)
        cnt.append(int(v))
        tot.append(mus*100/tasknum)
    _, ax = plt.subplots()
    plt.title("Greater 3h MUTIL-CARD WORKLOAD CDF ")
    ax.set_xlabel('MIN/MAX*100 ')
    ax.set_ylabel('Task Count')
    # ax.set_ylim(bottom=0)
    # cnt[0]=120
    # cnt[1]=110

    xx = np.arange(-1,101)
    ax.bar(xx,cnt)
    # plt.text(xx[0],cnt[0],1621,ha='center')
    # plt.xticks(rotation=45)
    # plt.text(xx[1],cnt[1],299,ha='center')
    # plt.xticks(rotation=45)
    ax2=ax.twinx()
    ax2.set_ylabel('Occupy/%')
    #ax2.set_ylim(bottom=0)
    ax2.plot(xx,tot,c="orange")

    plt.tight_layout()
    plt.savefig("./plot/cdf-Muti_Card.jpg")
    plt.show()
    #print(tot)
    file.close()

def Getgpucpu():
    file = open("./1.log", 'r', encoding='utf-8')
    str = file.read()
    data = str.split(" ")
    tot = []
    cnt = []
    mus = 0.0
    tasknum = 2079
    for v in data :
        mus = mus + float(v)
        cnt.append(int(v))
        tot.append(mus*100/tasknum)
        
    fig, ax = plt.subplots()
    plt.title("Greater 3h GPU-CPU CDF ")
    ax.set_xlabel('GPU-CPU ')
    ax.set_ylabel('Task Count')
    # ax.set_ylim(bottom=0)
    # cnt[100]=150
    # cnt[1]=110

    xx = np.arange(-100,101)
    print(len(xx))
    print(len(cnt))
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
    #plt.savefig("./plot/cdf-Muti_Card.jpg")
    plt.show()
    #print(tot)
    file.close()