import os
from subprocess import run, Popen, CREATE_NEW_CONSOLE, TimeoutExpired

import schedule
import time
from datetime import datetime,timedelta

def run_task():
    print("Do task...")

    year_n =3
    try_n = 1000

    timeout_s = 60 * 20
    failed_ls = []
    
    now = datetime.now()
    year_ls = []
    for i in range(year_n):
        year_ls.append(str(now.year - i))
    if now.month > 9:
        year_ls = [str(now.year + 1)] + year_ls

    for ls in ["LB","SC"]:
        for c in ["LIN","MSC","SRC","WAC","EAC","YSC"]:
            for y in year_ls:
                print(c,y,ls)

                cmd_str = "start /wait F:/Go_projects/scraper_two/main.exe "+c+" "+str(y)[-2:]+" "+ls+" "+str(try_n)
                subp = Popen(cmd_str, creationflags=CREATE_NEW_CONSOLE, shell=True, encoding="utf-8")
                try:
                    outs, errs = subp.communicate(timeout=timeout_s)
                except TimeoutExpired:
                    print(f'Timeout for "{cmd_str}" ({timeout_s}s) expired')
                    failed_ls.append([c,y,ls])
                    subp.kill()
                    outs, errs = subp.communicate()

    ## if failed, try until succeed
    while len(failed_ls) > 0:
        print("Re-try:",failed_ls)
        f_i = failed_ls[0]
        cmd_str = "start /wait C:/Users/hurui/OneDrive/temp/scraper_two/main.exe "+f_i[0]+" "+str(f_i[1])[-2:]+" "+f_i[2]+" "+str(try_n)
        subp = Popen(cmd_str, creationflags=CREATE_NEW_CONSOLE, shell=True, encoding="utf-8")
        try:
            outs, errs = subp.communicate(timeout=timeout_s)
            failed_ls.remove(f_i)
        except TimeoutExpired:
            print(f'Timeout for "{cmd_str}" ({timeout_s}s) expired')
            subp.kill()
            outs, errs = subp.communicate()

    print("Do task...Done!")


## Run job every day at specific HH:MM and next HH:MM:SS
# schedule.every().day.at("21:00:00").do(run_task)
#
# while True:
#     print("Running...",datetime.now())
#     schedule.run_pending()
#     time.sleep(1200)

run_task()






