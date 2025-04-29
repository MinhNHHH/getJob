"""
Config module of social listening
"""
import os
from dotenv import load_dotenv
from functools import partial

ENV_DEFAULT = {}

def find_dot_env_file():
	if os.environ.get("PYTHON_ENV", None) == "test":
		return ".env.test"
	else:
		return ".env"

def load_env():
	"""
	Loading all environments variable that has a SL_ prefix and strip it.
	SL_CRAWLER_CONNECTION_URI => CRAWLER_CONNECTION_URI
	To get this env, call `get_default("CRAWLER_CONNECTION_URI")`
	"""
	load_dotenv(find_dot_env_file())
	cfg = ENV_DEFAULT
	# TODO: add supports for nested config using env
	# Probably split by "."
	for env in os.environ:
		if env.startswith("GJ_"):
			env_name = env[len("GJ_"):]
			cfg[env_name] = os.environ[env]
	return cfg

def get(cfg, *path, required=True):
	"""
	Get a config with path.
	`path` could be a string for a top-level config or an array for nested config.

	config = {
	"first": {
		"second": "42"
	}
	}

	# get all level config
	>> get(config)
	=> {
	"first": {
		"second": "42"
	}
	}

	# get top level config
	>> get(config, "first")
	=> {
		"second": "42"
	}

	# get nested config
	>> get(config, "first", "second")
	=> 42

	Parameters
	----------

	path: str
		a path to config

	requried: boolean
		if `True`, raise an exception if the config is not found

	Returns
	-------
	dict | str | None
	The config for given path. In case config not found returns `None`
	"""
	def _get(cfg, k, required):
		try:
			cfg = cfg[k]
		except Exception:
			if required:
				raise Exception(f"Failed to get config with path: {path}")
				return None
		return cfg

	value = cfg
	for k in path:
		value = _get(value, k, required)
	return value

def get_int(cfg, *path, required=False):
	"""
	Like `get` but try to cast the returned values to integer.
	"""
	value = get(cfg, *path, required=required)
	return None if value == None else int(value)

########################### Default ###########################

CFG = load_env()
get_default = partial(get, CFG)          # get_default("jobs")
get_int_default = partial(get_int, CFG)  # get_int_default("jobs", "interval", "day")