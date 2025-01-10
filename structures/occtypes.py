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
                    "xvalues":[0,0],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 0,"sd": 0}},
                        {"type": "NormalDistribution","parameters":{"mean": 0,"sd": 0}}
                    ]
                }
            }
        }
    }

    return(output)

def build_damage_function(df):
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2024a)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[],
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

        output['damagefunctions']['depth']['damagefunction']['xvalues'].append(row['flood_depth'])
        output['damagefunctions']['depth']['damagefunction']['ydistributions'].append(ydist)

    return(output)

def build_mv_ghg_damage_function1():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2024a)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": 6103,
                    "depth": 484.8,
                    "sqft": 8.291,
                    "n_bed": 276.1,
                    "n_bath": -117.1,
                    "n_car": 718.1,
                    "depth_sqft": 0.7202,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0
                },
                "damagevectorsd": {
                    "intercept": 435.8,
                    "depth": 38.69,
                    "sqft": 2.352,
                    "n_bed": -18.22,
                    "n_bath": -15.26,
                    "n_car": 0,
                    "depth_sqft": 0,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0
                }
            }
        }
    }
    return(output)

def build_mv_ghg_damage_function2():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2024a)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": 6899,
                    "depth": 395,
                    "sqft": 8.456,
                    "n_bed": 285.5,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 0.7067,
                    "depth_n_bed": 0,
                    "depth_n_bath": -19.86,
                    "depth_n_car": 82.05
                },
                "damagevectorsd": {
                    "intercept": 171.3,
                    "depth": 26.01,
                    "sqft": 1.861,
                    "n_bed": 186.4,
                    "n_bath": -2.658,
                    "n_car": -155.7,
                    "depth_sqft": 0.01462,
                    "depth_n_bed": -4.091,
                    "depth_n_bath": 1.512,
                    "depth_n_car": 6.937
                }
            }
        }
    }
    return(output)

def build_mv_ghg_damage_function_null():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2024a)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": 0,
                    "depth": 0,
                    "sqft": 0,
                    "n_bed": 0,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 0,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0
                },
                "damagevectorsd": {
                    "intercept": 0,
                    "depth": 0,
                    "sqft": 0,
                    "n_bed": 0,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 0,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0
                }
            }
        }
    }
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
    dfs['co2_cost_pct_sd'] = (dfs['co2_cost_pct_mean'] - dfs['co2_cost_pct_low']) / 1.96
    dfs = dfs[['occtype', 'flood_depth', 'co2_cost_pct_mean', 'co2_cost_pct_sd']]
    dfs['flood_depth'] = dfs['flood_depth'].round(1)
    dfs = dfs[(dfs['flood_depth'] % 1 == 0) & (dfs['flood_depth'] <= 16)]
    print(dfs)


def main():

    dfs = pd.read_parquet("rowan_2024a_dmg_fns.parquet")
    with open("occtypes_original.json", "r") as f:
        occtypes = json.load(f)

    dfs['co2_cost_pct_sd'] = (dfs['co2_cost_pct_mean'] - dfs['co2_cost_pct_low']) / 1.96
    dfs = dfs[['occtype', 'flood_depth', 'co2_cost_pct_mean', 'co2_cost_pct_sd']]
    dfs['flood_depth'] = dfs['flood_depth'].round(1)
    dfs = dfs[(dfs['flood_depth']*10 % 1 == 0) & (dfs['flood_depth'] <= 16)]

    dfs1 = dfs[dfs['occtype'] == "RES1-1S"]
    dfs2 = dfs[dfs['occtype'] == "RES1-2S"]

    df1 = build_damage_function(dfs1)
    df2 = build_damage_function(dfs2)
    dfnull = build_null_df(dfs1)

    df_mv_ghg1 = build_mv_ghg_damage_function1()
    df_mv_ghg2 = build_mv_ghg_damage_function2()
    df_mv_ghg_null = build_mv_ghg_damage_function_null()
    
    occtypes_out = {"occupancytypes":{}}
    for key, o in occtypes['occupancytypes'].items():
        occtypes_out["occupancytypes"][key] = o
        if(o['name'][0:7] == "RES1-1S"):
            print(f"{o['name']} getting ghg df1")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = df1
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas2'] = df_mv_ghg1
        elif(o['name'][0:7] in ["RES1-2S", "RES1-3S", "RES1-SL", "RES3A", "RES3B"]):
            print(f"{o['name']} getting ghg df2")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = df2
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas2'] = df_mv_ghg2
        else:
            print(f"{o['name']} getting ghg dfnull")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = dfnull
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas2'] = df_mv_ghg_null

    # print(occtypes_out)
    # with open("occtypes_ghgrowan2024a.json", "w") as out:
    with open("occtypes.json", "w") as out:
        json.dump(occtypes_out, out, indent=4)




if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.realpath(__file__)))
    main()