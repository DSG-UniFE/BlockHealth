from psutil import virtual_memory, cpu_percent
from datetime import datetime
from time import sleep


with open("cpu-mem-values.txt", 'a') as file:
    file.write("Started CPU/memory usage\n\n")
    file.write("Date                | CPU(%) | MEM(%) | MEM (bytes)\n")
    while True:
        dt = datetime.now().strftime("%d/%m/%Y %H:%M:%S")
        cpu = cpu_percent()
        mem_perc = virtual_memory().percent
        mem_act = virtual_memory().active
        s = dt + " | " + str(cpu) + "% | " + str(mem_perc) + "% | " + str(mem_act) + "\n"
        file.write(s)
        sleep(1)


