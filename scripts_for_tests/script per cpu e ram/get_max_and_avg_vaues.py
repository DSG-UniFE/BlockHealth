from statistics import mean
from datetime import datetime
import sys

# if len(sys.argv) != 3:
#     print("usage: {} start stop".format(sys.argv[0]))
#     exit(-1)

# date and time examples
# '2019-09-06 16:30:04.989 UTC'
start = '2019-09-10 17:21:29.919 UTC'
stop = '2019-09-10 17:22:37.897 UTC'


try:
    datetime.strptime(start, "%Y-%m-%d %H:%M:%S.%f %Z")
except:
    print("incorrect START date format, usage: YYYY-MM-DD HH:MM:SS.FFF UTC")

try:
    datetime.strptime(stop, "%Y-%m-%d %H:%M:%S.%f %Z")
except:
    print("incorrect STOP date format, usage: YYYY-MM-DD HH:MM:SS.FFF UTC")

    
date_list = []
cpu_list = []
mem_list = []
start_time = datetime.strptime(start, "%Y-%m-%d %H:%M:%S.%f %Z")
stop_time = datetime.strptime(stop, "%Y-%m-%d %H:%M:%S.%f %Z")
with open("cpu-mem-values.txt") as file:
    line = file.read()
lines = line.split("\n")
for el in lines:
    if el == "\n":
        print("el: ", el)
        continue
    if len(el) > 0 and el[0] in '0123456789':
        date, cpu, mem, _ = el.split(" | ")
        # print("date: ", date)
        # print("cpu: ", cpu)
        # print("mem: ", mem)
        # 06/09/2019 16:35:59
        # "%Y-%m-%d %H:%M:%S.%f %Z"
        d = datetime.strptime(date, "%d/%m/%Y %H:%M:%S")
        if d > start_time and d < stop_time:
            cpu_list.append(float(cpu.strip("%")))
            mem_list.append(float(mem.strip("%")))
    
cpu_avg = mean(cpu_list)
mem_avg = mean(mem_list)
cpu_max = max(cpu_list)
mem_max = max(mem_list)

with open("cpu-mem-values.txt", 'a') as file:
    string = "\nCPU(%) avg value: {}%\n".format(round(cpu_avg, 1))
    file.write(string)
    print(string)
    string = "CPU(%) max value: {}%\n".format(round(cpu_max, 1))
    file.write(string)
    print(string)
    string = "MEM(%) avg value: {}%\n".format(round(mem_avg, 1))
    file.write(string)
    print(string)
    string = "MEM(%) max value: {}%\n".format(round(mem_max, 1))
    file.write(string)
    print(string)
