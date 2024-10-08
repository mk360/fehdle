/**
 * search for hero
 * on input, just show hero name and set id in the dom
 * on suggestion click, start comparison
 * horizontal display:
 * * 2-rows: character portrait
 * * 1 row: move type icon, weapon type icon, Game Name, Color,
 * Special properties (Ascended, Duo, Legendary, Duo Legendary, etc.)
 */

document.getElementById("guess-search").onchange = (e) => {
    if (e.target.value.trim().length >= 2) {
        fetch(`http://localhost:4444/hero?q=${encodeURIComponent(e.target.value)}`).then((res) => res.json()).then(console.log);
    }
}
