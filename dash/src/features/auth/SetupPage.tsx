import { SignupForm } from "./components";


export default function SetupPage() {
  return (
    <div className="grid min-h-svh w-screen max-h-screen overflow-clip lg:grid-cols-2">
      {/* Added justify-center and min-h-svh to center the content vertically */}
      <div className="flex flex-col z-20 gap-4 p-6 md:p-10 w-screen justify-center min-h-svh">
        <div className="flex justify-center gap-2 md:justify-start">
          <p className="flex items-center gap-2 font-medium text-xl" >
            <img src="/mist.png" alt="Mist Logo" className="size-8" />
            Mist
          </p>
        </div>
        <div className="flex flex-1 mx-auto items-center w-full justify-center">
          <div className="w-full max-w-md">
            <SignupForm />
          </div>
        </div>
      </div>
      <img
        src="/cloud-computing.png"
        alt="Image"
        className="absolute inset-0 h-full z-0 w-full object-cover dark:brightness-[0.3] dark:grayscale"
      />
    </div>
  )
}
