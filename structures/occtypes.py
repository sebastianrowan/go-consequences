# import numpy as np
# import pandas as pd
# import geopandas as gpd
import os
import json
from pprint import pprint


def prettify_occtypes():
    os.chdir(os.path.dirname(os.path.realpath(__file__)))
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)
    
    with open("occtypes_sbr.json", "w") as out:
        json.dump(occtypes, out, indent=4)

def add_damage_functions():
    # read occtypes json
    # for occtype in occtypes:
    #   if occtype == RES1-1S or RES1-2S:
    #       if damage driver == depth:
    #           insert damage function in correct place
    #   copy occtype to output json
    pass

def main():
    # prettify_occtypes()

    print("Done!")

if __name__ == "__main__":
    main()