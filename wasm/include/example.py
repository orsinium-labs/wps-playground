import sys
from sys import *
from typing import Dict, List, Union

assert True

number_with_useless_plus = +5
string_modifier = R'(\n)'
CONSTANT = []

leading_zero = 1.2e01
secondary_slice = items[1:][:3]

print('test')


def complex_annotation(
    first: List[Union[List[str], Dict[str, Dict[str, str]]]],
):
    ...


sum_container = 0  # re-implementing `sum`
for sum_item in file_obj:
    sum_container += sum_item
