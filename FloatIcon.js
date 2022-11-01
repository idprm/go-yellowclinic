import Link from 'next/link'

export default function FloatIcon() {
  return (
    <div className='fixed z-0 bottom-20 right-6'>
      <div className='flex-col justify-center items-center'>
        <Link href="https://wa.me/6281281881802?text=Layanan%20Homecare%20Perawatan%20Medis">
          <a target="_blank">
            <img src={`/images/icons/whatsapp.png`} className='p-1 w-12 h-12 drop-shadow-lg hover:animate-bounce duration-300' />
          </a>
        </Link>
        <Link href="https://t.me/CepatSehat">
          <a target="_blank">
            <img src={`/images/icons/telegram.png`} className='p-1 w-12 h-12 drop-shadow-lg hover:animate-bounce duration-300' />
          </a>
        </Link>
      </div>
    </div>
  )
}