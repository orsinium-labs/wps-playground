import json
import site
import sys
import os
from pathlib import Path
from random import choice
from string import ascii_lowercase
from textwrap import dedent
from unittest.mock import Mock, patch


# reset cache of packages (see importlib.metadata.FastPath.mtime)
os.utime(site.getsitepackages()[0])

# for type annotations in flake8
sys.modules['multiprocessing'] = Mock()
sys.modules['multiprocessing.pool'] = Mock()
patch("multiprocessing.pool", Mock())

if True:  # lazy way to avoid isort and autopep8 moving it
    from flake8.main.application import Application
    from flake8.formatting.base import BaseFormatter
    from flake8.style_guide import Violation

# set from the caller by patching globals
config: str
text: str


def random_name():
    name = ''
    for _ in range(20):
        name += choice(ascii_lowercase)
    return name + '.py'


# save flake8 config
Path("setup.cfg").write_text(config)  # noqa: F821

# save source code
path = Path(random_name())
path.write_text(dedent(text))  # noqa: F821


class Formatter(BaseFormatter):
    def after_init(self):
        self._out = []
        return super().after_init()

    def _write(self, output: str) -> None:
        self._out.append(output)

    def format(self, error: Violation) -> str:
        filename = error.filename
        if filename.startswith('./'):
            filename = filename[2:]
        return json.dumps(dict(
            path=filename,
            code=error.code,
            description=error.text,

            line=error.line_number,
            column=error.column_number,

            context=error.physical_line,
        ))


class App(Application):
    def make_formatter(self, formatter_class=None):
        self.formatter = Formatter(self.options)


# run flake8
app = App()
code = 0
try:
    app.run(['--jobs=0', str(path)])
except SystemExit as err:
    code = int(err.args[0])
else:
    code = app.exit_code()

# remove file
path.unlink()
