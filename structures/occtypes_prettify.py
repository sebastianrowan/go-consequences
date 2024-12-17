# import numpy as np
# import pandas as pd
# import geopandas as gpd
import os
import json
from pprint import pprint

def main():
    os.chdir(os.path.dirname(os.path.realpath(__file__)))
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)
    
    with open("occtypes_sbr.json", "w") as out:
        json.dump(occtypes, out, indent=4)

if __name__ == "__main__":
    main()