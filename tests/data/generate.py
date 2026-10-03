from pathlib import Path
from PIL import Image

root = Path(__file__).parent
rgb = Image.new('RGB', (8, 8), (32, 96, 160))
rgb.save(root / 'solid.jpg', quality=100, subsampling=0)
rgb.save(root / 'progressive.jpg', quality=100, progressive=True, subsampling=0)
gray = Image.new('L', (8, 8), 128)
gray.save(root / 'gray.webp', quality=100)
rgb.save(root / 'color.webp', quality=100)
Image.new('RGB', (8, 8), (0, 0, 0)).save(root / 'black.webp', quality=100)
Image.new('RGB', (8, 8), (255, 255, 255)).save(root / 'white.webp', quality=100)
Image.new('RGBA', (8, 8), (32, 96, 160, 128)).save(root / 'color-alpha.webp', quality=100)
rgba = Image.new('RGBA', (2, 2))
rgba.putdata([(255, 0, 0, 255), (0, 255, 0, 128), (0, 0, 255, 64), (255, 255, 255, 0)])
rgba.save(root / 'alpha.webp', lossless=True, exact=True)
(root.parents[1] / 'examples/basic/tests/data/alpha.webp').write_bytes((root / 'alpha.webp').read_bytes())
