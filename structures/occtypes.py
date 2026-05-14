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
                "source": "Rowan et al. (2025a)",
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
                "source": "Rowan et al. (2025a)",
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

def build_mv_dmg_damage_function1():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2025b)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": -3835,
                    "depth": 13391,
                    "sqft": -8.043,
                    "n_bed": 0,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 16.361,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.9208
                },
                "damagevectorsd": {
                    "intercept": -405.8,
                    "depth": 1232,
                    "sqft": 2.353,
                    "n_bed": 0,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 2.389,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.803
                }
            }
        }
    }
    return(output)

def build_mv_dmg_damage_function2():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2025b)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": -10090,
                    "depth": 15690,
                    "sqft": -6.864,
                    "n_bed": -2127,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 9.379,
                    "depth_n_bed": 5136,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.8921
                },
                "damagevectorsd": {
                    "intercept": -3074,
                    "depth": 2205,
                    "sqft": 0.9845,
                    "n_bed": -93.3,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 1.402,
                    "depth_n_bed": 1047,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.8085
                }
            }
        }
    }
    return(output)

def build_mv_ghg_damage_function1():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2025b)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": -1085,
                    "depth": 3388,
                    "sqft": -5.666,
                    "n_bed": 0,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 5.876,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.9382
                },
                "damagevectorsd": {
                    "intercept": -158,
                    "depth": 357.9,
                    "sqft": -1.14,
                    "n_bed": 0,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 1.03,
                    "depth_n_bed": 0,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.8449
                }
            }
        }
    }
    return(output)

def build_mv_ghg_damage_function2():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2025b)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": -1822,
                    "depth": 3556,
                    "sqft": -4.333,
                    "n_bed": -1048,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 3.492,
                    "depth_n_bed": 1543,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.9114
                },
                "damagevectorsd": {
                    "intercept": 388,
                    "depth": 103.3,
                    "sqft": -0.7937,
                    "n_bed": -363.6,
                    "n_bath": 0,
                    "n_car": 0,
                    "depth_sqft": 0.6293,
                    "depth_n_bed": 312.8,
                    "depth_n_bath": 0,
                    "depth_n_car": 0,
                    "r-squared": 0.8373
                }
            }
        }
    }
    return(output)

def build_mv_damage_function_null():
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Rowan et al. (2025b)",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[1,2],
                    "ydistributions": [
                        {"type": "NormalDistribution","parameters":{"mean": 10,"sd": 0.1}},
                        {"type": "NormalDistribution","parameters":{"mean": 20,"sd": 0.2}}
                    ]
                },
                "damagevectormean": {
                    "intercept": -9999,
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

    df_mv_dmg1 = build_mv_dmg_damage_function1()
    df_mv_dmg2 = build_mv_dmg_damage_function2()
    df_mv_ghg1 = build_mv_ghg_damage_function1()
    df_mv_ghg2 = build_mv_ghg_damage_function2()
    df_mv_null = build_mv_damage_function_null()
    
    occtypes_out = {"occupancytypes":{}}
    for key, o in occtypes['occupancytypes'].items():
        occtypes_out["occupancytypes"][key] = o
        if(o['name'][0:7] == "RES1-1S"):
            print(f"{o['name']} getting ghg df1")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['mv_structure'] = df_mv_dmg1
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = df1
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas2'] = df_mv_ghg1
        elif(o['name'][0:7] in ["RES1-2S", "RES1-3S", "RES1-SL", "RES3A", "RES3B"]):
            print(f"{o['name']} getting ghg df2")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['mv_structure'] = df_mv_dmg2
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = df2
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas2'] = df_mv_ghg2
        else:
            print(f"{o['name']} getting ghg dfnull")
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['mv_structure'] = df_mv_null
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas'] = dfnull
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['greenhouse_gas2'] = df_mv_null

    # print(occtypes_out)
    # with open("occtypes_ghgrowan2024a.json", "w") as out:
    with open("occtypes.json", "w") as out:
        json.dump(occtypes_out, out, indent=4)




if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.realpath(__file__)))
    main()