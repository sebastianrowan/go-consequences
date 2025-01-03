# import numpy as np
import pandas as pd
# import geopandas as gpd
import os
import json
from pprint import pprint


def prettify_occtypes():
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)
    
    with open("occtypes_sbr.json", "w") as out:
        json.dump(occtypes, out, indent=4)

def build_null_df(df):
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2024a)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[-2,-1,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16],
                    "ydistributions": []
                }
            }
        }
    }
    
    for i, row in df.iterrows():
        ydist = {
            "type": "NormalDistribution",
            "parameters":{
                "mean": 0,
                "sd": 0
            }
        }

        output['damagefunctions']['depth']['damagefunction']['ydistributions'].append(ydist)

    return(output)

def build_damage_function(df):
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2024a)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[-2,-1,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16],
                    "ydistributions": []
                }
            }
        }
    }
    
    for i, row in df.iterrows():
        # $190/1000kg = $0.19/kg
        mean = row['co2_cost_pct_mean']/0.190 # output will be kg CO2eq per total replacement value
        sd = row['co2_cost_pct_sd']/0.190

        ydist = {
            "type": "NormalDistribution",
            "parameters":{
                "mean": mean,
                "sd": sd
            }
        }

        output['damagefunctions']['depth']['damagefunction']['ydistributions'].append(ydist)

    return(output)


def read_occtypes():
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)
    
    for ot in occtypes['occupancytypes'].keys():
        print(ot)
        # print(occtypes['occupancytypes'][ot]['componentdamagefunctions']['structure']['damagefunctions'].keys())
        # print("----")
    
def print_dfs():
    dfs = pd.read_parquet("rowan_2024a_dmg_fns.parquet")
    dfs['co2_cost_pct_sd'] = (dfs['co2_cost_pct_mean'] - dfs['co2_cost_pct_low']) * 0.5
    dfs = dfs[['occtype', 'flood_depth', 'co2_cost_pct_mean', 'co2_cost_pct_sd']]
    dfs['flood_depth'] = dfs['flood_depth'].round(1)
    dfs = dfs[(dfs['flood_depth'] % 1 == 0) & (dfs['flood_depth'] <= 16)]
    print(dfs)


def main():

    dfs = pd.read_parquet("rowan_2024a_dmg_fns.parquet")
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)

    dfs['co2_cost_pct_sd'] = (dfs['co2_cost_pct_mean'] - dfs['co2_cost_pct_low']) * 0.5
    dfs = dfs[['occtype', 'flood_depth', 'co2_cost_pct_mean', 'co2_cost_pct_sd']]
    dfs['flood_depth'] = dfs['flood_depth'].round(1)
    dfs = dfs[(dfs['flood_depth'] % 1 == 0) & (dfs['flood_depth'] <= 16)]

    dfs1 = dfs[dfs['occtype'] == "RES1-1S"]
    dfs2 = dfs[dfs['occtype'] == "RES1-2S"]

    df1 = build_damage_function(dfs1)
    df2 = build_damage_function(dfs2)
    dfnull = build_null_df(dfs1)
    
    occtypes_out = {"occupancytypes":{}}
    for key, o in occtypes['occupancytypes'].items():
        occtypes_out["occupancytypes"][key] = o
        if(o['name'][0:7] == "RES1-1S"):
            print(f"{o['name']} getting ghg df1")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = df1
        elif(o['name'][0:7] in ["RES1-2S", "RES1-3S", "RES1-SL", "RES3A", "RES3B"]):
            print(f"{o['name']} getting ghg df2")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = df2
        else:
            print(f"{o['name']} getting ghg dfnull")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = dfnull

    # print(occtypes_out)
    # with open("occtypes_ghgrowan2024a.json", "w") as out:
    with open("occtypes.json", "w") as out:
        json.dump(occtypes_out, out, indent=4)




if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.realpath(__file__)))
    main()