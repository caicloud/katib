import grpc

from pkg.api.suggestion.v1alpha3.python import advisor_pb2_grpc
from pkg.api.suggestion.v1alpha3.python import advisor_pb2

DEFAULT_PORT="0.0.0.0:6789"

def run():
    channel = grpc.insecure_channel(DEFAULT_PORT)
    stub = advisor_pb2.AdvisorSuggestionStub(channel)
    experiment = advisor_pb2.Experiment(
        name="test",
        experiment_spec=advisor_pb2.ExperimentSpec(
            algorithm=advisor_pb2.AlgorithmSpec(
                algorithm_name="RandomSearch"
            ),
            objective=advisor_pb2.ObjectiveSpec(
                type=advisor_pb2.MAXIMIZE,
                goal=0.9
            ),
            parameter_specs=advisor_pb2.ParameterSpecs(
                parameters=[
                     advisor_pb2.ParameterSpec(
                        name="param1",
                        parameter_type=advisor_pb2.INT,
                        feasible_space=advisor_pb2.FeasibleSpace(max="5", min="1", list=[]),
                    ),
                    advisor_pb2.ParameterSpec(
                        name="param2",
                        parameter_type=advisor_pb2.CATEGORICAL,
                        feasible_space=advisor_pb2.FeasibleSpace(max=None, min=None, list=["cat1", "cat2", "cat3"])
                    ),
                    advisor_pb2.ParameterSpec(
                        name="param3",
                        parameter_type=advisor_pb2.DISCRETE,
                        feasible_space=advisor_pb2.FeasibleSpace(max=None, min=None, list=["3", "2", "6"])
                    ),
                    advisor_pb2.ParameterSpec(
                        name="param4",
                        parameter_type=advisor_pb2.DOUBLE,
                        feasible_space=advisor_pb2.FeasibleSpace(max="5", min="1", list=[])
                    )
                ]
            )
        )
    )
    
    response = stub.GetSuggestions(advisor_pb2.GetAdvisorSuggestionsRequest(
        experiment=experiment,
        trials=[advisor_pb2.Trial(name="test")],
        request_number=2,
    ))
    print(response)
    

if __name__ == "__main__":
    run()
