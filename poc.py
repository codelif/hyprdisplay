import curses
import curses.ascii
from curses.textpad import rectangle
from typing import List, Self


default_config = {
    "displays": {
        "mon-1": {
            "name": "eDP-1",
            "res": (1280, 720),
            "scale": 1,
            "pos": (0, 0),
            "primary": True,
        },
        "mon-2": {
            "name": "HDMI-1",
            "res": (1280, 720),
            "scale": 2 / 3,
            "pos": (1280, -1000),
            "primary": False,
        },
        "mon-3": {
            "name": "HDMI-2",
            "res": (3840, 2160),
            "scale": 2,
            "pos": (1280, 0),
            "primary": False,
        },
    },
    "primary": "mon-1",
}


def generate_monitors(config: dict):
    displays: dict = config["displays"]
    primary = displays[config["primary"]]

    displays_without_primary = displays.copy()
    displays_without_primary.pop(config["primary"])
    rest: list = list(displays_without_primary.values())

    pscale = primary["scale"]
    pres = [int(i / pscale) for i in primary["res"]]

    # Global Pixels per Cell ratio, calculated with primary display resolution as reference.
    global_ppc = (pres[0] // 24, pres[1] // 8)
    primary_lines = pres[1] // global_ppc[1]
    primary_cols = pres[0] // global_ppc[0]
    origin = ((curses.LINES - primary_lines) // 2, (curses.COLS - primary_cols) // 2)
    primary_monitor = Monitor(
        origin[0],
        origin[1],
        primary_lines,
        primary_cols,
        primary,
    )

    mon_objs = []
    for disp in rest:
        disp_lines = int((disp["res"][1] / disp["scale"]) / global_ppc[1])
        disp_cols = int((disp["res"][0] / disp["scale"]) / global_ppc[0])
        mon = Monitor(
            origin[0] + ((disp["pos"][1] - primary["pos"][1]) // global_ppc[1]),
            origin[1] + ((disp["pos"][0] - primary["pos"][0]) // global_ppc[0]),
            disp_lines,
            disp_cols,
            disp,
        )

        mon_objs.append(mon)

    return primary_monitor, mon_objs, global_ppc


def main(stdscr: curses.window):
    MOVE_KEYS = [
        curses.KEY_UP,
        curses.KEY_DOWN,
        curses.KEY_LEFT,
        curses.KEY_RIGHT,
        ord("k"),
        ord("j"),
        ord("h"),
        ord("l"),
    ]
    curses.set_escdelay(25)
    curses.curs_set(0)

    # Clear screen
    stdscr.clear()
    stdscr.refresh()
    curses.use_default_colors()

    mon1, wow, ppc = generate_monitors(default_config)
    refreshall(stdscr, mon1, *wow)
    while (ch := stdscr.getch()) not in [curses.ascii.ESC, ord("q")]:
        if ch in MOVE_KEYS:
            if ch in [curses.KEY_UP, ord("k")]:
                mon1.move_rel(-1, 0)
            elif ch in [curses.KEY_DOWN, ord("j")]:
                mon1.move_rel(1, 0)
            elif ch in [curses.KEY_LEFT, ord("h")]:
                mon1.move_rel(0, -2)
            elif ch in [curses.KEY_RIGHT, ord("l")]:
                mon1.move_rel(0, 2)

        stdscr.clear()
        status = f"({curses.LINES}, {curses.COLS}) ({mon1.posy}, {mon1.posx}) ({ppc[0]}, {ppc[1]})"
        stdscr.addstr(curses.LINES - 1, curses.COLS - len(status) - 1, status)
        # stdscr.addstr(0, 0, "Collision List:")
        # for i, mon in enumerate(wow, start=1):
        #     stdscr.addstr(i, 0, f"{mon.config['name']}: {mon1.is_colliding_with(mon)}")

        refreshall(stdscr, mon1, *wow)

    return calculate_config(mon1, wow, ppc)


def refreshall(*wins):
    for win in wins:
        win.noutrefresh()

    curses.doupdate()


class Monitor:
    def __init__(
        self, posy: int, posx: int, sizey: int, sizex: int, config: dict
    ) -> None:
        self.sizey = sizey
        self.sizex = sizex
        self.__set_pos(posy, posx)
        self.window = curses.newpad(sizey + 1, sizex + 1)
        self.config = config
        self.redraw()

    def refresh(self):
        self.noutrefresh()
        curses.doupdate()

    def noutrefresh(self):
        self.window.noutrefresh(
            0,
            0,
            self.posy,
            self.posx,
            self.posy + self.sizey - 1,
            self.posx + self.sizex - 1,
        )

    def remove(self):
        self.window.clear()

    def redraw(self):
        rectangle(self.window, 0, 0, self.sizey - 1, self.sizex - 1)
        self.add_text(self.config["name"])

    def add_text(self, text: str):
        line = (self.sizey - 1) // 2
        col = (self.sizex - len(text)) // 2
        self.window.addstr(line, col, text)

    def __bounded_pos(self, posy, posx):
        if posy < 0:
            posy = 0
        elif posy > curses.LINES - self.sizey:
            posy = curses.LINES - self.sizey

        if posx < 0:
            posx = 0
        elif posx > curses.COLS - self.sizex:
            posx = curses.COLS - self.sizex

        return posy, posx

    def __set_pos(self, posy: int, posx: int):
        self.posy = posy
        self.posx = posx

        self.miny = self.posy
        self.maxy = self.posy + self.sizey

        self.minx = self.posx
        self.maxx = self.posx + self.sizex

    def move_abs(self, posy, posx):
        self.__set_pos(*self.__bounded_pos(posy, posx))

    def move_rel(self, rely, relx):
        new_posy = self.posy + rely
        new_posx = self.posx + relx

        self.__set_pos(*self.__bounded_pos(new_posy, new_posx))

    def is_inside(self, y: int, x: int) -> bool:
        in_row_segment = self.posy <= y <= self.posy + self.sizey
        in_cols_segment = self.posx <= x <= self.posx + self.sizex

        return in_row_segment and in_cols_segment

    def is_colliding_with(self, mon: Self) -> bool:
        # Using AABB as all displays are rectanges and are only rotated in 90° (atleast according to hyprland)
        # https://en.wikipedia.org/wiki/Minimum_bounding_box#Axis-aligned_minimum_bounding_box
        return (
            (self.minx < mon.maxx)
            and (self.maxx > mon.minx)
            and (self.miny < mon.maxy)
            and (self.maxy > mon.miny)
        )


def calculate_config(primary: Monitor, rest: List[Monitor], ppc):
    template = "monitor = {name},preferred,{x},{y}"
    configs = []

    configs.append(template.format(name=primary.config["name"], x=0, y=0))

    for mon in rest:
        name = mon.config["name"]
        x = (mon.posx - primary.posx) * ppc[1]
        y = (mon.posy - primary.posy) * ppc[0]
        configs.append(template.format(name=name, x=x, y=y))

    return configs


if __name__ == "__main__":
    from pprint import pprint

    configs = curses.wrapper(main)

    print(*configs, sep="\n")
    # generate_monitors(default_conf g)
