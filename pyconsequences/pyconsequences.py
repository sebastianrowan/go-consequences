# import numpy as np
# import pandas as pd
# import geopandas as gpd
import os
import errno
import ctypes

#TODO: there should be logic here to check OS and import .so if user is on linux or mac
gclib = ctypes.CDLL("./pyconsequences.so")

_run_from_config_file = gclib.RunFromConfigFile
# _run_from_config_file.argtypes = [ctypes.c_char_p]

def run_from_config_file(filename: str) -> None:
    if (not os.path.exists(filename)):
        raise FileNotFoundError(
            errno.ENOENT, 
            os.strerror(errno.ENOENT), 
            filename
        )
    print(f"sending {filename} as config to go consequences")
    _run_from_config_file(filename.encode("utf-8"))
    
    


def main():
    run_from_config_file("config.json")

if __name__ == '__main__':
    main()