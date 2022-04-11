from zipfile import ZipFile
from base64 import b64decode
from io import BytesIO

archive: bytes
stream = BytesIO(b64decode(archive))  # noqa: F821
with ZipFile(stream) as zip_archive:
    zip_archive.extractall('.')
del archive  # noqa: F821
