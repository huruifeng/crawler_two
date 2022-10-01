import subprocess
import schedule
import time

def run_task():
    print("Do task...")
    try_n = 20
    for ls in ["LB","SC"]:
        for c in ["LIN","MSC","SRC","WAC","EAC","YSC"]:
            for y in [20,21,22]:
                print(c,y,ls,try_n)
                # Call your exe

                cmd_str = "C:/Users/hurui/OneDrive/temp/scraper_two/main.exe "+c+" "+str(y)+" "+ls+" "+str(try_n)
                subprocess.call(cmd_str,shell=True)

    # if you want to print output
    # p = subprocess.check_output('C:\pathtotool.exe -2 c:\data ')
    print("Do task...Done!")


# Run job every day at specific HH:MM and next HH:MM:SS
schedule.every().day.at("22:00:00").do(run_task)

while True:
    print("Running...")
    schedule.run_pending()
    time.sleep(1800)
