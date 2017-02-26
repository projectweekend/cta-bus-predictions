import json
import os


with open(os.path.expanduser('~/.cta_bus_predictions')) as cf:
    data = json.load(cf)
