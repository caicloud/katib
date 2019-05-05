from logging import getLogger, StreamHandler, INFO, DEBUG
import logging
import json

from pkg.suggestion.v1alpha3.advisor.algorithm.random_search import RandomSearchAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.grid_search import GridSearchAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.bayesian_optimization import BayesianOptimization
# from pkg.suggestion.v1alpha3.advisor.algorithm.tpe import TpeAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.simulate_anneal import SimulateAnnealAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.hyperopt_random_search import HyperoptRandomSearchAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.quasi_random_search import QuasiRandomSearchAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.chocolate_random_search import ChocolateRandomSearchAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.chocolate_grid_search import ChocolateGridSearchAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.chocolate_bayes import ChocolateBayesAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.cmaes import CmaesAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.mocmaes import MocmaesAlgorithm
# from pkg.suggestion.v1alpha3.advisor.algorithm.skopt_bayesian_optimization import SkoptBayesianOptimization

from pkg.api.suggestion.v1alpha3.python import advisor_pb2_grpc
from pkg.api.suggestion.v1alpha3.python import advisor_pb2


class AdvisorUnifiedService(advisor_pb2_grpc.AdvisorSuggestionServicer):
    def __init__(self, logger=None):
        if logger == None:
            self.logger = getLogger(__name__)
            FORMAT = '%(asctime)-15s %(message)s'
            logging.basicConfig(format=FORMAT)
            handler = StreamHandler()
            handler.setLevel(INFO)
            self.logger.setLevel(INFO)
            self.logger.addHandler(handler)
            self.logger.propagate = False

    def GetSuggestions(self, request, context):
        experiment = request.experiment
        trials = request.trials
        study_configuration = self.composeConfig(experiment)
        input_trials = self.composeTrials(trials)
        if experiment.experiment_spec.algorithm.algorithm_name == "RandomSearch":
            algorithm = RandomSearchAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "GridSearch":
        #     algorithm = GridSearchAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "BayesianOptimization":
        #     algorithm = BayesianOptimization()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "TPE":
        #     algorithm = TpeAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "HyperoptRandomSearch":
        #     algorithm = HyperoptRandomSearchAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "SimulateAnneal":
        #     algorithm = SimulateAnnealAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "QuasiRandomSearch":
        #     algorithm = QuasiRandomSearchAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "ChocolateRandomSearch":
        #     algorithm = ChocolateRandomSearchAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "ChocolateGridSearch":
        #     algorithm = ChocolateGridSearchAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "ChocolateBayes":
        #     algorithm = ChocolateBayesAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "CMAES":
        #     algorithm = CmaesAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "MOCMAES":
        #     algorithm = MocmaesAlgorithm()
        # elif experiment.experiment_spec.algorithm.algorithm_name == "SkoptBayesianOptimization":
        #     algorithm = SkoptBayesianOptimization()
        else:
            # TODO(gaocegege): Add error return.
            return advisor_pb2.GetAdvisorSuggestionsReply(
                trials=[]
            )
        res = algorithm.get_new_suggestions({"name": experiment.name},
                                            study_configuration, input_trials, request.request_number)
        return advisor_pb2.GetAdvisorSuggestionsReply(
            trials=self.decodeTrials(res)
        )

    def composeConfig(self, experiment):
        """
        compose the study configuration. The study's study_configuration is like this.
        {
            "goal": "MAXIMIZE",
            "params": [
                {
                    "parameterName": "hidden1",
                    "type": "INTEGER",
                    "minValue": 40,
                    "maxValue": 400,
                    "scalingType": "LINEAR"
                }
            ]
        }
        """
        config = {
            "params": []
        }
        if experiment.experiment_spec.objective.type == advisor_pb2.MAXIMIZE:
            config["goal"] = "MAXIMIZE"
        elif experiment.experiment_spec.objective.type == advisor_pb2.MINIMIZE:
            config["goal"] = "MINIMIZE"
        else:
            return config

        for p in experiment.experiment_spec.parameter_specs.parameters:
            # TODO(gaocegege): Error check.
            config["params"].append(getAdvisorParameter(p))
        return config

    def composeTrials(self, trials):
        """
        The trial's parameter_values_json should be like this.
        {
            "hidden1": 40
        }
        """
        return {}

    def decodeTrials(self, trials):
        res = []
        for t in trials:
            assignments = []
            values = json.loads(t["parameter_values"])
            for name, value in values.items():
                assignments.append(
                    advisor_pb2.ParameterAssignment(name=name, value=str(value)))
            rt = advisor_pb2.Trial(
                name=t["name"],
                spec=advisor_pb2.TrialSpec(
                    parameter_assignments=advisor_pb2.ParameterAssignments(
                        assignments=assignments
                    )
                )
            )
            res.append(rt)
        return res


def getAdvisorParameter(p):
    if p.parameter_type == advisor_pb2.INT:
        return {
            "parameterName": p.name,
            "type": "INTEGER",
            "minValue": float(p.feasible_space.min),
            "maxValue": float(p.feasible_space.max),
            "scalingType": "LINEAR",
        }
    elif p.parameter_type == advisor_pb2.DOUBLE:
        return {
            "parameterName": p.name,
            "type": "DOUBLE",
            "minValue": float(p.feasible_space.min),
            "maxValue": float(p.feasible_space.max),
            "scalingType": "LINEAR",
        }
    elif p.parameter_type == advisor_pb2.CATEGORICAL:
        fp = ','.join(p.feasible_space.list)
        return {
            "parameterName": p.name,
            "type": "CATEGORICAL",
            "feasiblePoints": fp,
            "scalingType": "LINEAR",
        }
    elif p.parameter_type == advisor_pb2.DISCRETE:
        fp = ''.join(p.feasible_space.list)
        return {
            "parameterName": p.name,
            "type": "DISCRETE",
            "feasiblePoints": fp,
            "scalingType": "LINEAR",
        }
    else:
        return {}
