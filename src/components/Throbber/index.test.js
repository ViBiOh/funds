import React from "react";
import { shallow } from "enzyme";
import Throbber from "./index";

it("should render into a div", () => {
  expect(shallow(<Throbber />).type()).toEqual("div");
});

it("should have no label by default", () => {
  expect(shallow(<Throbber />).find("span").length).toEqual(0);
});

it("should have label when given", () => {
  expect(
    shallow(<Throbber label="test" />)
      .find("span")
      .text()
  ).toEqual("test");
});
