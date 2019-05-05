import json
import random

from pkg.suggestion.v1alpha3.advisor.algorithm.abstract_algorithm import AbstractSuggestionAlgorithm
from pkg.suggestion.v1alpha3.advisor.algorithm.util import AlgorithmUtil


class RandomSearchAlgorithm(AbstractSuggestionAlgorithm):
  def get_new_suggestions(self, study, study_configuration_json, trials=[], number=1):
    """
    Get the new suggested trials with random search.
    """

    return_trial_list = []

    params = study_configuration_json["params"]

    for i in range(number):
      trial = {}
      trial["name"] = study["name"]

      parameter_values_json = {}

      for param in params:

        if param["type"] == "DOUBLE":
          suggest_value = AlgorithmUtil.get_random_value(
              param["minValue"], param["maxValue"])

        elif param["type"] == "INTEGER":
          suggest_value = AlgorithmUtil.get_random_int_value(
              param["minValue"], param["maxValue"])

        elif param["type"] == "DISCRETE":
          feasible_point_list = [
              float(value.strip())
              for value in param["feasiblePoints"].split(",")
          ]
          suggest_value = AlgorithmUtil.get_random_item_from_list(
              feasible_point_list)

        elif param["type"] == "CATEGORICAL":
          feasible_point_list = [
              value.strip() for value in param["feasiblePoints"].split(",")
          ]
          suggest_value = AlgorithmUtil.get_random_item_from_list(
              feasible_point_list)

        parameter_values_json[param["parameterName"]] = suggest_value

      trial["parameter_values"] = json.dumps(parameter_values_json)
      return_trial_list.append(trial)

    return return_trial_list
