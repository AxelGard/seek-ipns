from . import random_forest
from . import score

def main(model_name) -> dict:
    if model_name == "randomforest":
        return tune_random_forest()

    return {}



def tune_random_forest() -> dict:
    result = {}

    # default
    score.setup()
    model = random_forest.RandomForest()
    result["default"] = score.get_scores(model)

    # Clanging data 
    model = random_forest.RandomForest(clean=False)
    result["raw"] = score.get_scores(model)

    # nr estimators 
    for n in range(440,480, 5):
        model = random_forest.RandomForest(n_estimators=n)
        #result[f"{n}-estimator"] = score.get_scores(model)["RandomForest"]
    
    # nr estimators 
    for n in [nr / 10.0 for nr in range(4,7)]:
        model = random_forest.RandomForest(prob_min=n)
        result[f"{n}-prob"] = score.get_scores(model)

    # best estimation 
    model = random_forest.RandomForest(clean=False, n_estimators=450, prob_min=0.4)
    result["best_guess"] = score.get_scores(model)

    return result
