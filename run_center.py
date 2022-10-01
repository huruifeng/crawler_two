import os
import subprocess

import schedule
import time
from datetime import datetime,timedelta

def run_task():
    print("Do task...")

    now = datetime.now()
    year_n =3
    year_ls = []
    for i in range(year_n):
        year_ls.append(str(now.year - i))
    if now.month > 9:
        year_ls = [str(now.year + 1)] + year_ls

    for ls in ["LB","SC"]:
        for c in ["LIN","MSC","SRC","WAC","EAC","YSC"]:
            for y in year_ls:
                print(c,y,ls)
                # Call your exe
                cmd_str = "start /wait F:/Go_projects/scraper_two/main.exe "+c+" "+str(y)[-2:]+" "+ls
                os.system(cmd_str)

    # if you want to print output
    # p = subprocess.check_output('C:\pathtotool.exe -2 c:\data ')
    print("Do task...Done!")


# # Run job every day at specific HH:MM and next HH:MM:SS
# schedule.every().day.at("22:00:00").do(run_task)
#
# while True:
#     print("Running...")
#     schedule.run_pending()
#     time.sleep(1800)

run_task()
