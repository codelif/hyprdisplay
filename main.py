import curses
import curses.ascii
from curses.textpad import rectangle
from typing import Tuple


default_config = {
    "displays": {
        "mon-1": {
            "name": "eDP-1",
            "res": (1920, 1080),
            "scale": 1,
            "pos": (0, 0),
            "primary": True,
        },
        "mon-2": {
            "name": "HDMI-1",
            "res": (1280, 720),
            "scale": 2 / 3,
            "pos": (0, -1080),
            "primary": False,
        },
        "mon-3": {
            "name": "HDMI-2",
            "res": (3840, 2160),
            "scale": 1,
            "pos": (1920, -500),
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
    global_ppc = (pres[0] // 30, pres[1] // 10)
    primary_lines = pres[1] // global_ppc[1]
    primary_cols = pres[0] // global_ppc[0]
    origin = ((curses.LINES - primary_lines) // 2, (curses.COLS - primary_cols) // 2)
    primary_monitor = Monitor(
        origin[0],
        origin[1],
        primary_lines,
        primary_cols,
    )

    primary_monitor.refresh()

    mon_objs = []
    for disp in rest:
        disp_lines = int((disp["res"][1] / disp["scale"]) / global_ppc[1])
        disp_cols = int((disp["res"][0] / disp["scale"]) / global_ppc[0])
        mon = Monitor(
            origin[0] + (disp["pos"][1] // global_ppc[1]),
            origin[1] + (disp["pos"][0] // global_ppc[0]),
            disp_lines,
            disp_cols,
        )

        mon.refresh()
        mon_objs.append(mon)

    return primary_monitor, mon_objs


def main(stdscr: curses.window):
    curses.set_escdelay(25)
    curses.curs_set(0)

    # Clear screen
    stdscr.clear()
    stdscr.refresh()
    curses.use_default_colors()

    mon1, wow = generate_monitors(default_config)

    # MON_LINES = 10
    # MON_COLS = 30
    # mon1 = Monitor(
    #     (curses.LINES - MON_LINES) // 2,
    #     (curses.COLS - MON_COLS) // 2,
    #     MON_LINES,
    #     MON_COLS,
    # )
    mon1.refresh()
    while (ch := stdscr.getch()) not in [curses.ascii.ESC, ord("q")]:
        if ch in [
            curses.KEY_UP,
            curses.KEY_DOWN,
            curses.KEY_LEFT,
            curses.KEY_RIGHT,
            ord("k"),
            ord("j"),
            ord("h"),
            ord("l"),
        ]:
            if ch in [curses.KEY_UP, ord("k")]:
                mon1.move_rel(-1, 0)
            elif ch in [curses.KEY_DOWN, ord("j")]:
                mon1.move_rel(1, 0)
            elif ch in [curses.KEY_LEFT, ord("h")]:
                mon1.move_rel(0, -2)
            elif ch in [curses.KEY_RIGHT, ord("l")]:
                mon1.move_rel(0, 2)

        stdscr.clear()
        status = f"({curses.LINES}, {curses.COLS}) ({mon1.posy}, {mon1.posx}) ({mon1.sizey}, {mon1.sizex})"
        stdscr.addstr(curses.LINES - 1, curses.COLS - len(status) - 1, status)
        refreshall(stdscr, mon1, *wow)

    stdscr.getch()


def refreshall(*wins):
    for win in wins:
        win.noutrefresh()

    curses.doupdate()


class Monitor:
    def __init__(self, posy, posx, sizey, sizex) -> None:
        self.posy = posy
        self.posx = posx
        self.sizey = sizey
        self.sizex = sizex
        self.window = curses.newpad(sizey + 1, sizex + 1)
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

    def move_abs(self, posy, posx):
        self.posy, self.posx = self.__bounded_pos(posy, posx)

    def move_rel(self, rely, relx):
        new_posy = self.posy + rely
        new_posx = self.posx + relx

        self.posy, self.posx = self.__bounded_pos(new_posy, new_posx)


if __name__ == "__main__":
    from pprint import pprint

    curses.wrapper(main)
    # generate_monitors(default_conf g)
