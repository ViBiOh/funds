import funtch from "funtch";
import ConfigService from "./index";

jest.mock("funtch");
jest.mock("AppStore");

it("should fetch config from env", () => {
  funtch.get.mockResolvedValue({});

  return ConfigService.getConfig().then(() => {
    expect(funtch.get).toHaveBeenCalledWith("/env");
  });
});
